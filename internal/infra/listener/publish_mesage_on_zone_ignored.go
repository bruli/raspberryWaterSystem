package listener

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PublishMessageOnZoneIgnored struct {
	ch     cqs.CommandHandler
	tracer trace.Tracer
}

func (p PublishMessageOnZoneIgnored) Listen(ctx context.Context, ev cqs.Event) error {
	ctx, span := p.tracer.Start(ctx, "PublishMessageOnZoneIgnored.Listen")
	defer span.End()
	event, _ := ev.(zone.Ignored)
	message := spew.Sprintf("Zone %q ignored to send water. Reason: %s", event.ZoneName, event.Reason)
	if _, err := p.ch.Handle(ctx, app.PublishMessageCmd{Message: message}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "message published")
	return nil
}

func NewPublishMessageOnZoneIgnored(ch cqs.CommandHandler, tracer trace.Tracer) *PublishMessageOnZoneIgnored {
	return &PublishMessageOnZoneIgnored{ch: ch, tracer: tracer}
}
