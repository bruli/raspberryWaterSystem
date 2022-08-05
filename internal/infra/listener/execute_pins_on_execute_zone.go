package listener

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
)

type ExecutePinsOnExecuteZone struct {
	ch cqs.CommandHandler
}

func (e ExecutePinsOnExecuteZone) Listen(ctx context.Context, ev cqs.Event) error {
	event, _ := ev.(zone.Executed)
	_, err := e.ch.Handle(ctx, app.ExecutePinsCmd{
		Seconds: event.Seconds,
		Pins:    event.RelayPins,
	})
	return err
}

func NewExecutePinsOnExecuteZone(ch cqs.CommandHandler) ExecutePinsOnExecuteZone {
	return ExecutePinsOnExecuteZone{ch: ch}
}
