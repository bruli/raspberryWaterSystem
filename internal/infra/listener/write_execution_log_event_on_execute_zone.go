package listener

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type WriteExecutionLogEventOnExecuteZone struct {
	eventsRepo *disk.EventsRepository
}

func (w WriteExecutionLogEventOnExecuteZone) Listen(ctx context.Context, ev cqs.Event) error {
	event, _ := ev.(zone.Executed)
	evnt, err := disk.NewFromExecutionLog(ctx, &disk.Log{
		Seconds:    int(event.Seconds),
		ZoneName:   event.ZoneName,
		ExecutedAt: event.EventAt(),
	})
	if err != nil {
		return err
	}
	return w.eventsRepo.Save(ctx, evnt)
}

func NewWriteExecutionLogEventOnExecuteZone(eventsRepo *disk.EventsRepository) *WriteExecutionLogEventOnExecuteZone {
	return &WriteExecutionLogEventOnExecuteZone{eventsRepo: eventsRepo}
}
