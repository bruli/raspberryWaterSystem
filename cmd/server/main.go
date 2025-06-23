package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/config"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	infrahttp "github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/internal/infra/listener"
	"github.com/bruli/raspberryWaterSystem/internal/infra/memory"
	"github.com/bruli/raspberryWaterSystem/internal/infra/telegram"
	"github.com/bruli/raspberryWaterSystem/internal/infra/worker"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/rs/zerolog"
)

func main() {
	log := buildLogger()
	conf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed building config")
	}
	ctx := context.Background()

	eventsCh := make(chan cqs.Event, 5)
	defer close(eventsCh)

	logCHMdw := cqs.NewCommandHndErrorMiddleware(&log)
	eventsCHMdw := app.NewEventMiddleware(eventsCh)
	eventsMultiCHMdw := cqs.CommandHandlerMultiMiddleware(logCHMdw, eventsCHMdw)
	logQHMdw := cqs.NewQueryHndErrorMiddleware(&log)

	tr := temperatureRepository()
	rr := rainRepository()
	sr := memory.NewStatusRepository()
	zr := disk.NewZoneRepository(conf.ZonesFile)
	dailyRepo := disk.NewProgramRepository(conf.DailyProgramsFile)
	oddRepo := disk.NewProgramRepository(conf.OddProgramsFile)
	evenRepo := disk.NewProgramRepository(conf.EvenProgramsFile)
	weeklyRepo := disk.NewWeeklyRepository(conf.WeeklyProgramsFile)
	tempProgRepo := disk.NewTemperatureProgramRepository(conf.TemperatureProgramsFile)
	execLogRepo := disk.NewExecutionLogRepository(conf.ExecutionLogsFile)
	pe := pinsExecutor()
	messagePublisher := telegram.NewMessagePublisher(conf.TelegramToken, conf.TelegramChatID)

	qhBus := app.NewQueryBus()
	qhBus.Subscribe(app.FindWeatherQueryName, logQHMdw(app.NewFindWeather(tr, rr)))
	qhBus.Subscribe(app.FindStatusQueryName, logQHMdw(app.NewFindStatus(sr)))
	qhBus.Subscribe(app.FindAllProgramsQueryName, logQHMdw(app.NewFindAllPrograms(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	qhBus.Subscribe(app.FindProgramsInTimeQueryName, logQHMdw(app.NewFindProgramsInTime(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	qhBus.Subscribe(app.FindExecutionLogsQueryName, logQHMdw(app.NewFindExecutionLogs(execLogRepo)))
	qhBus.Subscribe(app.FindZonesQueryName, logQHMdw(app.NewFindZones(zr)))

	chBus := app.NewCommandBus()
	chBus.Subscribe(app.CreateStatusCmdName, logCHMdw(app.NewCreateStatus(sr)))
	chBus.Subscribe(app.UpdateStatusCmdName, logCHMdw(app.NewUpdateStatus(sr)))
	chBus.Subscribe(app.CreateZoneCmdName, logCHMdw(app.NewCreateZone(zr)))
	chBus.Subscribe(app.ExecuteZoneCmdName, eventsMultiCHMdw(app.NewExecuteZone(zr)))
	chBus.Subscribe(app.ExecutePinsCmdName, logCHMdw(app.NewExecutePins(pe)))
	chBus.Subscribe(app.SaveExecutionLogCmdName, logCHMdw(app.NewSaveExecutionLog(execLogRepo)))
	chBus.Subscribe(app.PublishMessageCmdName, logCHMdw(app.NewPublishMessage(messagePublisher)))
	chBus.Subscribe(app.RemoveZoneCmdName, logCHMdw(app.NewRemoveZone(zr)))
	chBus.Subscribe(app.ActivateDeactivateServerCmdName, logCHMdw(app.NewActivateDeactivateServer(sr)))
	chBus.Subscribe(app.ExecuteZoneWithStatusCmdName, eventsMultiCHMdw(app.NewExecuteZoneWithStatus(zr, sr)))
	chBus.Subscribe(app.UpdateZoneCommandName, logCHMdw(app.NewUpdateZone(zr)))
	chBus.Subscribe(app.CreateDailyProgramCommandName, logCHMdw(app.NewCreateDailyProgram(dailyRepo, zr)))
	chBus.Subscribe(app.CreateOddProgramCommandName, logCHMdw(app.NewCreateOddProgram(oddRepo, zr)))
	chBus.Subscribe(app.CreateEvenProgramCommandName, logCHMdw(app.NewCreateEvenProgram(evenRepo, zr)))
	chBus.Subscribe(app.RemoveDailyProgramCommandName, logCHMdw(app.NewRemoveDailyProgram(dailyRepo)))
	chBus.Subscribe(app.RemoveOddProgramCommandName, logCHMdw(app.NewRemoveOddProgram(oddRepo)))
	chBus.Subscribe(app.RemoveEvenProgramCommandName, logCHMdw(app.NewRemoveEvenProgram(evenRepo)))
	chBus.Subscribe(app.CreateWeeklyProgramCommandName, logCHMdw(app.NewCreateWeeklyProgram(weeklyRepo)))
	chBus.Subscribe(app.RemoveWeeklyProgramCommandName, logCHMdw(app.NewRemoveWeeklyProgram(weeklyRepo)))
	chBus.Subscribe(app.CreateTemperatureProgramCommandName, logCHMdw(app.NewCreateTemperatureProgram(tempProgRepo)))
	chBus.Subscribe(app.RemoveTemperatureProgramCommandName, logCHMdw(app.NewRemoveTemperatureProgram(tempProgRepo)))
	chBus.Subscribe(app.UpdateTemperatureProgramCommandName, logCHMdw(app.NewUpdateTemperatureProgram(tempProgRepo)))

	eventBus := cqs.NewEventBus()
	eventBus.Subscribe(zone.Executed{
		BasicEvent: cqs.BasicEvent{NameAttr: zone.ExecutedEventName},
	}, listener.NewExecutePinsOnExecuteZone(chBus))
	eventBus.Subscribe(zone.Ignored{
		BasicEvent: cqs.BasicEvent{NameAttr: zone.IgnoredEventName},
	}, listener.NewPublishMessageOnZoneIgnored(chBus))

	if err = initStatus(ctx, chBus, qhBus); err != nil {
		log.Fatal().Err(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go updateStatusWorker(ctx, qhBus, chBus, &log)
	go eventsWorker(ctx, eventsCh, eventBus, &log)
	go executionInTimeWorker(ctx, qhBus, chBus, &log)

	go runTelegramBot(conf, qhBus, chBus, log, ctx)
	runHTTPServer(chBus, qhBus, conf, ctx, log)
}

func runTelegramBot(conf *config.Config, qhBus app.QueryBus, chBus app.CommandBus, log zerolog.Logger, ctx context.Context) {
	if !conf.TelegramBotEnabled {
		log.Info().Msg("[TELEGRAM SERVICE] disabled")
		return
	}
	telegramServer, err := telegram.NewCommandReader(conf.TelegramToken, qhBus, chBus)
	if err != nil {
		log.Fatal().Err(err).Msgf("[TELEGRAM SERVICE] failed building telegram server: %s", err)
	}
	telegramServer.Read(ctx, &log)
}

func runHTTPServer(chBus app.CommandBus, qhBus app.QueryBus, conf *config.Config, ctx context.Context, log zerolog.Logger) {
	definitions := handlersDefinition(chBus, qhBus, conf.AuthToken)
	httpHandlers := infrahttp.NewHandler(definitions)
	if err := infrahttp.RunServer(ctx, conf.ServerURL, httpHandlers, &infrahttp.CORSOpt{}, &log); err != nil {
		log.Fatal().Err(err).Msg("system error")
	}
}

func buildLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return log
}

func executionInTimeWorker(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, logger *zerolog.Logger) {
	logger.Info().Msg("[WORKER] Execution in time: started")
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("[WORKER] Execution in time: context done")
			return
		case <-ticker.C:
			if err := worker.ExecutionInTime(ctx, qh, ch, vo.TimeNow()); err != nil {
				logger.Err(err).Msg("[WORKER] failed execution in time")
			}
		}
	}
}

func eventsWorker(ctx context.Context, ch <-chan cqs.Event, evBus cqs.EventBus, logger *zerolog.Logger) {
	logger.Info().Msg("[WORKER] Events: started")
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("[WORKER] Events: context done")
			return
		case event := <-ch:
			if err := evBus.Dispatch(ctx, event); err != nil {
				logger.Err(err).Msg(fmt.Sprintf("[WORKER] Events: failed dispatching %q", event.EventName()))
			}
		}
	}
}

func updateStatusWorker(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, logger *zerolog.Logger) {
	logger.Info().Msg("[WORKER] Update status: started")
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("[WORKER] Update status: context done")
			return
		case <-ticker.C:
			updateStatus(ctx, qh, ch, logger)
		}
	}
}

func updateStatus(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, logger *zerolog.Logger) {
	result, err := qh.Handle(ctx, app.FindWeatherQuery{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed getting weather")
	}
	weath, _ := result.(weather.Weather)
	if _, err = ch.Handle(ctx, app.UpdateStatusCmd{Weather: weath}); err != nil {
		logger.Fatal().Err(err).Msg("failed updating status")
	}
}

func initStatus(ctx context.Context, ch cqs.CommandHandler, qh cqs.QueryHandler) error {
	result, err := qh.Handle(ctx, app.FindWeatherQuery{})
	if err != nil {
		return err
	}
	weath, _ := result.(weather.Weather)
	_, err = ch.Handle(ctx, app.CreateStatusCmd{
		StartedAt: vo.TimeNow(),
		Weather:   weath,
	})
	return err
}

func handlersDefinition(chBus app.CommandBus, qhBus app.QueryBus, authToken string) infrahttp.HandlersDefinition {
	authMdw := infrahttp.AuthMiddleware(authToken)
	return infrahttp.HandlersDefinition{
		{
			Endpoint:    "/",
			Method:      http.MethodGet,
			HandlerFunc: infrahttp.Homepage(),
		},
		{
			Endpoint:    "/zones",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateZone(chBus)),
		},
		{
			Endpoint:    "/zones",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(infrahttp.FinZones(qhBus)),
		},
		{
			Endpoint:    "/zones/{id}/execute",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.ExecuteZone(chBus)),
		},
		{
			Endpoint:    "/zones/{id}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveZone(chBus)),
		},
		{
			Endpoint:    "/zones/{id}",
			Method:      http.MethodPut,
			HandlerFunc: authMdw(infrahttp.UpdateZone(chBus)),
		},
		{
			Endpoint:    "/status",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(infrahttp.FindStatus(qhBus)),
		},
		{
			Endpoint:    "/status/{action}",
			Method:      http.MethodPatch,
			HandlerFunc: authMdw(infrahttp.ActivateDeactivateServer(chBus)),
		},
		{
			Endpoint:    "/weather",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(infrahttp.FindWeather(qhBus)),
		},
		{
			Endpoint:    "/programs",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(infrahttp.FindAllPrograms(qhBus)),
		},
		{
			Endpoint:    "/programs/daily",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateProgram(chBus, infrahttp.DailyProgram)),
		},
		{
			Endpoint:    "/programs/daily/{hour}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveProgram(chBus, infrahttp.DailyProgram)),
		},
		{
			Endpoint:    "/programs/odd",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateProgram(chBus, infrahttp.OddProgram)),
		},
		{
			Endpoint:    "/programs/odd/{hour}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveProgram(chBus, infrahttp.OddProgram)),
		},
		{
			Endpoint:    "/programs/even",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateProgram(chBus, infrahttp.EvenProgram)),
		},
		{
			Endpoint:    "/programs/even/{hour}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveProgram(chBus, infrahttp.EvenProgram)),
		},
		{
			Endpoint:    "/programs/weekly",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateWeeklyProgram(chBus)),
		},
		{
			Endpoint:    "/programs/weekly/{day}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveWeeklyProgram(chBus)),
		},
		{
			Endpoint:    "/programs/temperature",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(infrahttp.CreateTemperatureProgram(chBus)),
		},
		{
			Endpoint:    "/programs/temperature/{temperature}",
			Method:      http.MethodPut,
			HandlerFunc: authMdw(infrahttp.UpdateTemperatureProgram(chBus)),
		},
		{
			Endpoint:    "/programs/temperature/{temperature}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(infrahttp.RemoveTemperatureProgram(chBus)),
		},
		{
			Endpoint:    "/logs",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(infrahttp.FindExecutionLogs(qhBus)),
		},
		{
			Endpoint:    "/metrics",
			Method:      http.MethodGet,
			HandlerFunc: infrahttp.Metrics(),
		},
	}
}
