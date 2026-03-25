package telegram

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type runnerCommand interface {
	CommandName() CommandName
}

type RunnerFunc func(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error

type runner interface {
	Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error
}

type runnerBus struct {
	runners map[CommandName]runner
	tracer  trace.Tracer
}

func (r *runnerBus) subscribe(command CommandName, run runner) {
	r.runners[command] = run
}

func (r *runnerBus) handle(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	ctx, span := r.tracer.Start(ctx, "runnerBus.handle")
	defer span.End()
	run, ok := r.runners[cmd.CommandName()]
	if !ok {
		err := errors.New("command not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if err := run.Run(ctx, chatID, msgs, cmd); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "command executed")
	return nil
}

func newRunnerBus(tracer trace.Tracer) *runnerBus {
	return &runnerBus{runners: make(map[CommandName]runner), tracer: tracer}
}
