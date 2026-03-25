package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const RemoveZoneCmdName = "removeZone"

type RemoveZoneCmd struct {
	ID string
}

func (r RemoveZoneCmd) Name() string {
	return RemoveZoneCmdName
}

type RemoveZone struct {
	zr     ZoneRepository
	tracer trace.Tracer
}

func (r RemoveZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	ctx, span := r.tracer.Start(ctx, "RemoveZoneCmd")
	defer span.End()
	co, ok := cmd.(RemoveZoneCmd)
	if !ok {
		err := cqs.NewInvalidCommandError(RemoveZoneCmdName, cmd.Name())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	zo, err := r.zr.FindByID(ctx, co.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if err = r.zr.Remove(ctx, zo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetStatus(codes.Ok, "zone removed")
	return nil, nil
}

func NewRemoveZone(zr ZoneRepository, tracer trace.Tracer) RemoveZone {
	return RemoveZone{zr: zr, tracer: tracer}
}
