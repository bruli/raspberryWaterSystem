package listener

import (
	"context"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
)

type ExecutePinsOnExecuteZone struct {
	ch cqs.CommandHandler
}

func (e ExecutePinsOnExecuteZone) Listen(ctx context.Context, ev cqs.Event) error {
	now := time.Now()
	event, _ := ev.(zone.Executed)
	if _, err := e.ch.Handle(ctx, app.ExecutePinsCmd{
		Seconds: event.Seconds,
		Pins:    event.RelayPins,
	}); err != nil {
		return err
	}
	sec, _ := program.ParseSeconds(int(event.Seconds))
	_, err := e.ch.Handle(ctx, app.SaveExecutionLogCmd{
		ZoneName:   event.ZoneName,
		Seconds:    sec,
		ExecutedAt: now,
	})
	return err
}

func NewExecutePinsOnExecuteZone(ch cqs.CommandHandler) ExecutePinsOnExecuteZone {
	return ExecutePinsOnExecuteZone{ch: ch}
}
