package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

//go:generate moq -out zmock_execution_log_publisher_test.go --pkg app_test . ExecutionLogPublisher

type ExecutionLogPublisher interface {
	Publish(ctx context.Context, execLog program.ExecutionLog) error
}
