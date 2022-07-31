package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	logCHMdw := cqs.NewCommandHndErrorMiddleware(logger)
	logQHMdw := cqs.NewQueryHndErrorMiddleware(logger)

	tr := temperatureRepository()
	rr := rainRepository()
	sr := memory.StatusRepository{}

	qhBus := app.NewQueryBus()
	qhBus.Subscribe(app.FindWeatherQueryName, logQHMdw(app.NewFindWeather(tr, rr)))
	chBus := app.NewCommandBus()
	chBus.Subscribe(app.CreateStatusCmdName, logCHMdw(app.NewCreateStatus(sr)))

	if err = initStatus(ctx, chBus, qhBus); err != nil {
		log.Fatalln(err)
	}
	definitions, err := handlersDefinition(logger, chBus, logCHMdw, conf.ZonesFile(), conf.AuthToken())
	if err != nil {
		log.Fatalln(err)
	}
	httpHandlers := httpx.NewHandler(definitions)
	if err := httpx.RunServer(ctx, conf.ServerURL(), httpHandlers, &httpx.CORSOpt{}); err != nil {
		log.Fatalln(fmt.Errorf("system error: %w", err))
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

func handlersDefinition(log *log.Logger, chBus app.CommandBus, chMddw cqs.CommandHandlerMiddleware, zonesFile, authToken string) (httpx.HandlersDefinition, error) {
	zr := disk.NewZoneRepository(zonesFile)
	chBus.Subscribe(app.CreateZoneCmdName, chMddw(app.NewCreateZone(zr)))
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
	}, nil
}
