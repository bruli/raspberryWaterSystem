package listener

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/davecgh/go-spew/spew"
)

type PublishMessageOnZoneIgnored struct {
	ch cqs.CommandHandler
}

func (p PublishMessageOnZoneIgnored) Listen(ctx context.Context, ev cqs.Event) error {
	event, _ := ev.(zone.Ignored)
	message := spew.Sprintf("Zone %q ignored to send water.", event.ZoneName)
	if _, err := p.ch.Handle(ctx, app.PublishMessageCmd{Message: message}); err != nil {
		return err
	}
	return nil
}

func NewPublishMessageOnZoneIgnored(ch cqs.CommandHandler) *PublishMessageOnZoneIgnored {
	return &PublishMessageOnZoneIgnored{ch: ch}
}
