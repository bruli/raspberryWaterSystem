package server

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	"time"
)

type weatherDaemon struct {
	log  logger.Logger
	writ *weather.Writer
}

func newWeatherDaemon(log logger.Logger, writ *weather.Writer) *weatherDaemon {
	return &weatherDaemon{log: log, writ: writ}
}

func (w *weatherDaemon) execute(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.writ.Write()
			if err != nil {
				w.log.Fatalf("failed writting weather: %s", err)
			}
		}
	}
}
