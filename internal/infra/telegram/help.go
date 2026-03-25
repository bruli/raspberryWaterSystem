package telegram

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type helCommand struct{}

func (h helCommand) CommandName() CommandName {
	return HelpCommandName
}

type helpRunner struct {
	tracer trace.Tracer
}

func (h helpRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		_, span := h.tracer.Start(ctx, "helpRunner.Run")
		defer span.End()
		buildMessage(chatID, msgs, "Available commands:")
		for _, c := range initCommands() {
			buildMessage(chatID, msgs, fmt.Sprintf("%s -> %s, %q", c.name.String(), c.syntax, c.description))
		}
		span.SetStatus(codes.Ok, "help message sent")
		return nil
	}
}

func newHelpRunner(tracer trace.Tracer) *helpRunner {
	return &helpRunner{
		tracer: tracer,
	}
}
