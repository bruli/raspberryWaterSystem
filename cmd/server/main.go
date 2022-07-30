package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	logger := log.New(os.Stdout, config.ProjectPrefix, int(time.Now().Unix()))
	definitions, err := handlersDefinition(logger, conf.ZonesFile())
	if err != nil {
		log.Fatalln(err)
	}
	httpHandlers := httpx.NewHandler(definitions)
	if err := httpx.RunServer(ctx, conf.ServerURL(), httpHandlers, &httpx.CORSOpt{}); err != nil {
		log.Fatalln(fmt.Errorf("system error: %w", err))
	}
}

func handlersDefinition(log *log.Logger, zonesFile string) (httpx.HandlersDefinition, error) {
	zr := disk.NewZoneRepository(zonesFile)
	logCHMdw := cqs.NewCommandHndErrorMiddleware(log)
	chBus := app.NewCommandBus()
	chBus.Subscribe(app.CreateZoneCmdName, logCHMdw(app.NewCreateZone(zr)))
	return httpx.HandlersDefinition{
		{
			Endpoint:    "/",
			Method:      http.MethodGet,
			HandlerFunc: http2.Homepage(),
		},
		{
			Endpoint:    "/zones",
			Method:      http.MethodPost,
			HandlerFunc: http2.CreateZone(chBus),
		},
	}, nil
}
