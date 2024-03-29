package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/config"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/internal/infra/api"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/bruli/raspberryWaterSystem/internal/infra/fake"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/internal/infra/listener"
	"github.com/bruli/raspberryWaterSystem/internal/infra/memory"
	"github.com/bruli/raspberryWaterSystem/internal/infra/telegram"
	"github.com/bruli/raspberryWaterSystem/internal/infra/worker"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	conf, err := config.NewConfig()
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
	rr := rainRepository(conf)
	sr := memory.StatusRepository{}
	zr := disk.NewZoneRepository(conf.ZonesFile())
	dailyRepo := disk.NewProgramRepository(conf.DailyProgramsFile())
	oddRepo := disk.NewProgramRepository(conf.OddProgramsFile())
	evenRepo := disk.NewProgramRepository(conf.EvenProgramsFile())
	weeklyRepo := disk.NewWeeklyRepository(conf.WeeklyProgramsFile())
	tempProgRepo := disk.NewTemperatureProgramRepository(conf.TemperatureProgramsFile())
	execLogRepo := disk.NewExecutionLogRepository(conf.ExecutionLogsFile())
	pe := pinsExecutor()
	messagePublisher := telegram.NewMessagePublisher(conf.TelegramToken(), conf.TelegramChatID())

	qhBus := app.NewQueryBus()
	qhBus.Subscribe(app.FindWeatherQueryName, logQHMdw(app.NewFindWeather(tr, rr)))
	qhBus.Subscribe(app.FindStatusQueryName, logQHMdw(app.NewFindStatus(sr)))
	qhBus.Subscribe(app.FindAllProgramsQueryName, logQHMdw(app.NewFindAllPrograms(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	qhBus.Subscribe(app.FindProgramsInTimeQueryName, logQHMdw(app.NewFindProgramsInTime(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	qhBus.Subscribe(app.FindExecutionLogsQueryName, logQHMdw(app.NewFindExecutionLogs(execLogRepo)))

	chBus := app.NewCommandBus()
	chBus.Subscribe(app.CreateStatusCmdName, logCHMdw(app.NewCreateStatus(sr)))
	chBus.Subscribe(app.UpdateStatusCmdName, logCHMdw(app.NewUpdateStatus(sr)))
	chBus.Subscribe(app.CreateZoneCmdName, logCHMdw(app.NewCreateZone(zr)))
	chBus.Subscribe(app.CreateProgramsCmdName, logCHMdw(app.NewCreatePrograms(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	chBus.Subscribe(app.ExecuteZoneCmdName, eventsMultiCHMdw(app.NewExecuteZone(zr)))
	chBus.Subscribe(app.ExecutePinsCmdName, logCHMdw(app.NewExecutePins(pe)))
	chBus.Subscribe(app.SaveExecutionLogCmdName, logCHMdw(app.NewSaveExecutionLog(execLogRepo)))
	chBus.Subscribe(app.PublishMessageCmdName, logCHMdw(app.NewPublishMessage(messagePublisher)))
	chBus.Subscribe(app.RemoveZoneCmdName, logCHMdw(app.NewRemoveZone(zr)))
	chBus.Subscribe(app.ActivateDeactivateServerCmdName, logCHMdw(app.NewActivateDeactivateServer(sr)))
	chBus.Subscribe(app.ExecuteZoneWithStatusCmdName, eventsMultiCHMdw(app.NewExecuteZoneWithStatus(zr, sr)))

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

	go updateStatusWorker(ctx, qhBus, chBus)
	go eventsWorker(ctx, eventsCh, eventBus, &log)
	go executionInTimeWorker(ctx, qhBus, chBus, &log)

	definitions, err := handlersDefinition(chBus, qhBus, conf.AuthToken())
	if err != nil {
		log.Fatal().Err(err)
	}
	httpHandlers := httpx.NewHandler(definitions)
	if err = httpx.RunServer(ctx, conf.ServerURL(), httpHandlers, &httpx.CORSOpt{}); err != nil {
		log.Fatal().Err(err).Msg("system error")
	}
}

func executionInTimeWorker(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, logger *zerolog.Logger) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("execution in time worker context done")
			return
		case <-ticker.C:
			if err := worker.ExecutionInTime(ctx, qh, ch, vo.TimeNow()); err != nil {
				logger.Err(err).Msg("failed execution in time worker")
			}
		}
	}
}

func eventsWorker(ctx context.Context, ch <-chan cqs.Event, evBus cqs.EventBus, logger *zerolog.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			if err := evBus.Dispatch(ctx, event); err != nil {
				logger.Err(err).Msg(fmt.Sprintf("failed dispatching %q", event.EventName()))
			}
		}
	}
}

func rainRepository(conf config.Config) app.RainRepository {
	var rr app.RainRepository
	rr = fake.RainRepository{}
	if conf.Environment().IsProduction() {
		rr = api.NewRainRepository(conf.RainServerURL())
	}
	return rr
}

func updateStatusWorker(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("update status worker context done")
			return
		case <-ticker.C:
			updateStatus(ctx, qh, ch)
		}
	}
}

func updateStatus(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler) {
	result, err := qh.Handle(ctx, app.FindWeatherQuery{})
	if err != nil {
		log.Fatalln(err)
	}
	weath, _ := result.(weather.Weather)
	if _, err = ch.Handle(ctx, app.UpdateStatusCmd{Weather: weath}); err != nil {
		log.Fatalln(err)
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

func handlersDefinition(chBus app.CommandBus, qhBus app.QueryBus, authToken string) (httpx.HandlersDefinition, error) {
	authMdw := http2.AuthMiddleware(authToken)
	return httpx.HandlersDefinition{
		{
			Endpoint:    "/",
			Method:      http.MethodGet,
			HandlerFunc: http2.Homepage(),
		},
		{
			Endpoint:    "/zones",
			Method:      http.MethodPut,
			HandlerFunc: authMdw(http2.CreateZone(chBus)),
		},
		{
			Endpoint:    "/zones/{id}/execute",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(http2.ExecuteZone(chBus)),
		},
		{
			Endpoint:    "/zones/{id}",
			Method:      http.MethodDelete,
			HandlerFunc: authMdw(http2.RemoveZone(chBus)),
		},
		{
			Endpoint:    "/status",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(http2.FindStatus(qhBus)),
		},
		{
			Endpoint:    "/status/{action}",
			Method:      http.MethodPatch,
			HandlerFunc: authMdw(http2.ActivateDeactivateServer(chBus)),
		},
		{
			Endpoint:    "/weather",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(http2.FindWeather(qhBus)),
		},
		{
			Endpoint:    "/programs",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(http2.FindAllPrograms(qhBus)),
		},
		{
			Endpoint:    "/programs",
			Method:      http.MethodPut,
			HandlerFunc: authMdw(http2.CreatePrograms(chBus)),
		},
		{
			Endpoint:    "/logs",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(http2.FindExecutionLogs(qhBus)),
		},
	}, nil
}
