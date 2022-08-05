package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/infra/worker"

	"github.com/bruli/raspberryWaterSystem/internal/infra/listener"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"

	"github.com/bruli/raspberryWaterSystem/internal/infra/api"
	"github.com/bruli/raspberryWaterSystem/internal/infra/fake"

	"github.com/bruli/raspberryWaterSystem/internal/infra/memory"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/config"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	logger := log.New(os.Stdout, config.ProjectPrefix, int(time.Now().Unix()))

	eventsCh := make(chan cqs.Event)
	defer close(eventsCh)

	logCHMdw := cqs.NewCommandHndErrorMiddleware(logger)
	eventsCHMdw := app.NewEventMiddleware(eventsCh)
	eventsMultiCHMdw := cqs.CommandHandlerMultiMiddleware(logCHMdw, eventsCHMdw)
	logQHMdw := cqs.NewQueryHndErrorMiddleware(logger)

	tr := temperatureRepository()
	rr := rainRepository(conf)
	sr := memory.StatusRepository{}
	zr := disk.NewZoneRepository(conf.ZonesFile())
	dailyRepo := disk.NewProgramRepository(conf.DailyProgramsFile())
	oddRepo := disk.NewProgramRepository(conf.OddProgramsFile())
	evenRepo := disk.NewProgramRepository(conf.EvenProgramsFile())
	weeklyRepo := disk.NewWeeklyRepository(conf.WeeklyProgramsFile())
	tempProgRepo := disk.NewTemperatureProgramRepository(conf.TemperatureProgramsFile())
	pe := pinsExecutor(logger)

	qhBus := app.NewQueryBus()
	qhBus.Subscribe(app.FindWeatherQueryName, logQHMdw(app.NewFindWeather(tr, rr)))
	qhBus.Subscribe(app.FindStatusQueryName, logQHMdw(app.NewFindStatus(sr)))
	qhBus.Subscribe(app.FindAllProgramsQueryName, logQHMdw(app.NewFindAllPrograms(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	qhBus.Subscribe(app.FindProgramsInTimeQueryName, logQHMdw(app.NewFindProgramsInTime(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))

	chBus := app.NewCommandBus()
	chBus.Subscribe(app.CreateStatusCmdName, logCHMdw(app.NewCreateStatus(sr)))
	chBus.Subscribe(app.UpdateStatusCmdName, logCHMdw(app.NewUpdateStatus(sr)))
	chBus.Subscribe(app.CreateZoneCmdName, logCHMdw(app.NewCreateZone(zr)))
	chBus.Subscribe(app.CreateProgramsCmdName, logCHMdw(app.NewCreatePrograms(dailyRepo, oddRepo, evenRepo, weeklyRepo, tempProgRepo)))
	chBus.Subscribe(app.ExecuteZoneCmdName, eventsMultiCHMdw(app.NewExecuteZone(zr)))
	chBus.Subscribe(app.ExecutePinsCmdName, logCHMdw(app.NewExecutePins(pe)))

	eventBus := cqs.NewEventBus()
	eventBus.Subscribe(zone.Executed{
		BasicEvent: cqs.BasicEvent{NameAttr: zone.ExecutedEventName},
	}, listener.NewExecutePinsOnExecuteZone(chBus))

	if err = initStatus(ctx, chBus, qhBus); err != nil {
		log.Fatalln(err)
	}

	go updateStatusWorker(ctx, qhBus, chBus)
	go eventsWorker(ctx, eventsCh, eventBus, logger)
	go executionInTimeWorker(ctx, qhBus, chBus, logger)

	definitions, err := handlersDefinition(chBus, qhBus, conf.AuthToken())
	if err != nil {
		log.Fatalln(err)
	}
	httpHandlers := httpx.NewHandler(definitions)
	if err := httpx.RunServer(ctx, conf.ServerURL(), httpHandlers, &httpx.CORSOpt{}); err != nil {
		log.Fatalln(fmt.Errorf("system error: %w", err))
	}
}

func executionInTimeWorker(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, logger *log.Logger) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			logger.Println("execution in time worker context done")
			return
		case <-ticker.C:
			if err := worker.ExecutionInTime(ctx, qh, ch, time.Now()); err != nil {
				log.Printf("failed execution in time worker: %s", err.Error())
			}
		}
	}
}

func eventsWorker(ctx context.Context, ch <-chan cqs.Event, evBus cqs.EventBus, logger *log.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			if err := evBus.Dispatch(ctx, event); err != nil {
				logger.Printf("failed dispatching %s: %s", event.EventName(), err.Error())
			}
		}
	}
}

func rainRepository(conf config.Config) app.RainRepository {
	var rr app.RainRepository
	rr = fake.RainRepository{}
	if conf.Environment().IsProduction() {
		rr = api.RainRepository{}
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
		StartedAt: time.Now(),
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
			Method:      http.MethodPost,
			HandlerFunc: authMdw(http2.CreateZone(chBus)),
		},
		{
			Endpoint:    "/zones/{id}/execute",
			Method:      http.MethodPost,
			HandlerFunc: authMdw(http2.ExecuteZone(chBus)),
		},
		{
			Endpoint:    "/status",
			Method:      http.MethodGet,
			HandlerFunc: authMdw(http2.FindStatus(qhBus)),
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
			Method:      http.MethodPost,
			HandlerFunc: authMdw(http2.CreatePrograms(chBus)),
		},
	}, nil
}
