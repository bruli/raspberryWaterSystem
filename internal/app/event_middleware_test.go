package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestNewEventMiddleware(t *testing.T) {
	t.Run(`Given an event command handler middleware,
	when is built with an event channel,
	and command handler return events,
	then all events are write in channel`, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		errTest := errors.New("")
		event1 := Event{
			ID:        uuid.New(),
			Name:      "eventito 1",
			At:        time.Now(),
			AggRootID: "dkdkdkd",
		}
		eventCh := make(chan cqs.Event)
		eventMdw := app.NewEventMiddleware(eventCh)
		handler := eventMdw(commandHandler{
			events: []cqs.Event{
				event1,
			},
			err: errTest,
		})
		go func() {
			for {
				select {
				case <-ctx.Done():
					cancel()
					return
				case event := <-eventCh:
					require.Equal(t, event1, event)
				}
			}
		}()
		evnts, err := handler.Handle(ctx, command{})
		require.Nil(t, evnts)
		require.Equal(t, errTest, err)
	})
}

type Event struct {
	ID        uuid.UUID
	Name      string
	At        time.Time
	AggRootID string
}

func (e Event) EventID() uuid.UUID {
	return e.ID
}

func (e Event) EventName() string {
	return e.Name
}

func (e Event) EventAt() time.Time {
	return e.At
}

func (e Event) AggregateRootID() string {
	return e.AggRootID
}
