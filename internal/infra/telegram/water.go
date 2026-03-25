package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type waterCommand struct {
	arguments string
}

func (w waterCommand) CommandName() CommandName {
	return WaterCommandName
}

type waterRunner struct {
	ch     cqs.CommandHandler
	tracer trace.Tracer
}

func (w waterRunner) Run(ctx context.Context, _ int64, _ *Messages, cmd runnerCommand) error {
	ctx, span := w.tracer.Start(ctx, "waterRunner.Run")
	defer span.End()
	co, _ := cmd.(waterCommand)
	arguments := strings.Fields(co.arguments)
	if len(arguments) == 0 {
		err := fmt.Errorf("invalid arguments")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	zone := arguments[0]
	seconds, err := strconv.Atoi(arguments[1])
	if err != nil {
		err = fmt.Errorf("invalid seconds")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if _, err = w.ch.Handle(ctx, app.ExecuteZoneCmd{
		Seconds: uint(seconds),
		ZoneID:  zone,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to execute zone: %w", err)
	}
	span.SetStatus(codes.Ok, "zone executed")
	return nil
}

func newWaterRunner(ch cqs.CommandHandler, tracer trace.Tracer) *waterRunner {
	return &waterRunner{ch: ch, tracer: tracer}
}
