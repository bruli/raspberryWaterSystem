package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ActivateCommand struct{}

func (a ActivateCommand) CommandName() CommandName {
	return ActivateCommandName
}

type ActivateRunner struct {
	ch     cqs.CommandHandler
	tracer trace.Tracer
}

func (a ActivateRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	ctx, span := a.tracer.Start(ctx, "ActivateRunner.Run")
	defer span.End()
	if _, err := a.ch.Handle(ctx, app.ActivateDeactivateServerCmd{Active: true}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to activate: %w", err)
	}

	buildMessage(chatID, msgs, "Activated!!")
	span.SetStatus(codes.Ok, "activated")
	return nil
}

func NewActivateRunner(ch cqs.CommandHandler, tracer trace.Tracer) *ActivateRunner {
	return &ActivateRunner{ch: ch, tracer: tracer}
}
