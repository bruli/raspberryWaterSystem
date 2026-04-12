package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const ExecutePinsCmdName = "executePins"

type ExecutePinsCmd struct {
	Seconds uint
	Pins    []string
}

func (e ExecutePinsCmd) Name() string {
	return ExecutePinsCmdName
}

type ExecutePins struct {
	pe     PinExecutor
	tracer trace.Tracer
}

func (e ExecutePins) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	ctx, span := e.tracer.Start(ctx, "ExecutePinsCmd")
	defer span.End()
	co, ok := cmd.(ExecutePinsCmd)
	if !ok {
		err := cqs.NewInvalidCommandError(ExecutePinsCmdName, cmd.Name())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if err := e.pe.Execute(ctx, co.Seconds, co.Pins); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return nil, nil
}

func NewExecutePins(pe PinExecutor, tracer trace.Tracer) ExecutePins {
	return ExecutePins{pe: pe, tracer: tracer}
}
