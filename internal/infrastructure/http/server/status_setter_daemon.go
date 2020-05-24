package server

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	"time"
)

type statusSetterDaemon struct {
	log    logger.Logger
	setter *weather.StatusSetter
}

func newStatusSetterDaemon(log logger.Logger, setter *weather.StatusSetter) *statusSetterDaemon {
	return &statusSetterDaemon{log: log, setter: setter}
}

func (s *statusSetterDaemon) execute(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.setter.Set()
			if err != nil {
				s.log.Fatalf("status setter daemon failed: %s", err)
			}
		}
	}
}
