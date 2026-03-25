package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type logCommand struct {
	Number int
}

func (l logCommand) CommandName() CommandName {
	return LogCommandName
}

type logRunner struct {
	qh     cqs.QueryHandler
	tracer trace.Tracer
}

func (l logRunner) Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	ctx, span := l.tracer.Start(ctx, "logRunner.Run")
	defer span.End()
	co, _ := cmd.(logCommand)
	result, err := l.qh.Handle(ctx, app.FindExecutionLogsQuery{Limit: co.Number})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to find logs: %w", err)
	}
	logs, _ := result.([]string)
	if len(logs) == 0 {
		buildMessage(chatID, msgs, "No logs found")
		return nil
	}
	for _, lo := range logs {
		buildMessage(chatID, msgs, lo)
	}
	span.SetStatus(codes.Ok, "logs found")
	return nil
}

func newLogRunner(qh cqs.QueryHandler, tracer trace.Tracer) *logRunner {
	return &logRunner{qh: qh, tracer: tracer}
}
