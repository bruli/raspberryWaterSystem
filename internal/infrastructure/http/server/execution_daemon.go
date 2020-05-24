package server

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
)

type executionDaemon struct {
	log      logger.Logger
	executor *execution.Executor
}

func newExecutionDaemon(log logger.Logger, executor *execution.Executor) *executionDaemon {
	return &executionDaemon{log: log, executor: executor}
}

func (e *executionDaemon) execute(ctx context.Context, ch chan executionData) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			if err := e.executor.Execute(data.seconds, data.zone); err != nil {
				e.log.Fatalf("failed to execute execution: %s", err)
			}
		}
	}
}
