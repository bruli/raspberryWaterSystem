package cqs

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

// UnknownEventToDispatchError is self-described
type UnknownEventToDispatchError struct {
	event string
}

func (u UnknownEventToDispatchError) Error() string {
	return fmt.Sprintf("event %q is not declared to dispatch", u.event)
}

// EventListener is self-described
type EventListener interface {
	Listen(ctx context.Context, ev Event) error
}

// EventBus subscribe events with event listeners, and dispatch them.
type EventBus map[string][]EventListener

// NewEventBus is a constructor
func NewEventBus() EventBus {
	return make(map[string][]EventListener)
}

// Subscribe map event to event listeners
func (e EventBus) Subscribe(ev Event, listeners ...EventListener) {
	e[ev.EventName()] = listeners
}

// Dispatch execute event listeners from event
func (e EventBus) Dispatch(ctx context.Context, ev Event) error {
	mErr := vo.NewMultiError()
	list, ok := e[ev.EventName()]
	if !ok {
		return UnknownEventToDispatchError{event: ev.EventName()}
	}
	for _, l := range list {
		if err := l.Listen(ctx, ev); err != nil {
			mErr.Add(err)
		}
	}
	if mErr.HasErrors() {
		return mErr
	}
	return nil
}
