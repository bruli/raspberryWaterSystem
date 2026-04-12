package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type statusCommand struct{}

func (s statusCommand) CommandName() CommandName {
	return StatusCommandName
}

type statusRunner struct {
	qh     cqs.QueryHandler
	tracer trace.Tracer
}

func (s statusRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	ctx, span := s.tracer.Start(ctx, "statusRunner.Run")
	defer span.End()
	result, err := s.qh.Handle(ctx, app.FindStatusQuery{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to find status: %w", err)
	}
	st, _ := result.(status.Status)
	buildMessage(chatID, msgs, fmt.Sprintf("System started at: %s", st.SystemStartedAt().Format("2006-01-02 15:04:05")))
	buildMessage(chatID, msgs, fmt.Sprintf("Current temperature: %v *C", st.Weather().Temperature()))
	buildMessage(chatID, msgs, fmt.Sprintf("Current humidity: %v", st.Weather().Humidity()))
	buildMessage(chatID, msgs, fmt.Sprintf("Is raining: %v", st.Weather().IsRaining()))
	buildMessage(chatID, msgs, fmt.Sprintf("Active: %v", st.IsActive()))
	if st.UpdatedAt() != nil {
		buildMessage(chatID, msgs, fmt.Sprintf("System updated at: %s", st.UpdatedAt().Format("2006-01-02 15:04:05")))
	}
	span.SetStatus(codes.Ok, "status found")
	return nil
}

func newStatusRunner(qh cqs.QueryHandler, tracer trace.Tracer) *statusRunner {
	return &statusRunner{qh: qh, tracer: tracer}
}
