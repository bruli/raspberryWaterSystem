package server

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"time"
)

type executionInTimeDaemon struct {
	exec *execution.ExecutorInTime
	log  logger.Logger
}

func newExecutionInTimeDaemon(exec *execution.ExecutorInTime, log logger.Logger) *executionInTimeDaemon {
	return &executionInTimeDaemon{exec: exec, log: log}
}

func (e *executionInTimeDaemon) execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(60 * time.Second):
			e.log.Debug("reading execution in time...")
			err := e.exec.Execute(time.Now())
			if err != nil {
				e.log.Fatalf("failed to execute in time: %s", err)
			}
		}
	}
}
