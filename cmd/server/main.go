package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	definitions, err := handlersDefinition(logger)
	if err != nil {
		log.Fatalln(err)
	}
	httpHandlers := httpx.NewHandler(definitions)
	if err := httpx.RunServer(ctx, conf.ServerURL(), httpHandlers, &httpx.CORSOpt{}); err != nil {
		log.Fatalln(fmt.Errorf("system error: %w", err))
	}
}

func handlersDefinition(log *log.Logger) (httpx.HandlersDefinition, error) {
	//qhErrMdw := cqs.NewQueryHndErrorMiddleware(log)
	return httpx.HandlersDefinition{
		{
			Endpoint:    "/",
			Method:      http.MethodGet,
			HandlerFunc: http2.Homepage(),
		},
	}, nil
}
