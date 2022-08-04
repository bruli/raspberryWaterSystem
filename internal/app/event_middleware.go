package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

func NewEventMiddleware(evCh chan<- cqs.Event) cqs.CommandHandlerMiddleware {
	return func(h cqs.CommandHandler) cqs.CommandHandler {
		return cqs.CommandHandlerFunc(func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
			events, err := h.Handle(ctx, cmd)
			for _, ev := range events {
				evCh <- ev
			}
			return nil, err
		})
	}
}
