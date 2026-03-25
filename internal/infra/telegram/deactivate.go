package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type DeactivateCommand struct{}

func (a DeactivateCommand) CommandName() CommandName {
	return DeactivateCommandName
}

type DeactivateRunner struct {
	ch     cqs.CommandHandler
	tracer trace.Tracer
}

func (a DeactivateRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	ctx, span := a.tracer.Start(ctx, "DeactivateRunner.Run")
	defer span.End()
	if _, err := a.ch.Handle(ctx, app.ActivateDeactivateServerCmd{Active: false}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to deaactivate: %w", err)
	}

	buildMessage(chatID, msgs, "Deactivated!!")
	span.SetStatus(codes.Ok, "deactivated")
	return nil
}

func NewDeactivateRunner(ch cqs.CommandHandler, tracer trace.Tracer) *DeactivateRunner {
	return &DeactivateRunner{ch: ch, tracer: tracer}
}
