package listener

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type ExecutePinsOnExecuteZone struct {
	ch cqs.CommandHandler
}

func (e ExecutePinsOnExecuteZone) Listen(ctx context.Context, ev cqs.Event) error {
	now := vo.TimeNow()
	event, _ := ev.(zone.Executed)
	if _, err := e.ch.Handle(ctx, app.ExecutePinsCmd{
		Seconds: event.Seconds,
		Pins:    event.RelayPins,
	}); err != nil {
		return err
	}
	sec, _ := program.ParseSeconds(int(event.Seconds))
	if _, err := e.ch.Handle(ctx, app.SaveExecutionLogCmd{
		ZoneName:   event.ZoneName,
		Seconds:    sec,
		ExecutedAt: now,
	}); err != nil {
		return err
	}
	message := fmt.Sprintf("%s zone executed during %vs", event.ZoneName, sec.Int())
	_, err := e.ch.Handle(ctx, app.PublishMessageCmd{Message: message})
	return err
}

func NewExecutePinsOnExecuteZone(ch cqs.CommandHandler) ExecutePinsOnExecuteZone {
	return ExecutePinsOnExecuteZone{ch: ch}
}
