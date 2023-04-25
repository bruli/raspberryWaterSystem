package fake

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type PinsExecutor struct{}

func NewPinsExecutor() PinsExecutor {
	return PinsExecutor{}
}

func (p PinsExecutor) Execute(ctx context.Context, seconds uint, pins []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(time.Duration(seconds) * time.Second)
		log.Debug().Msgf("pins %s executed %v seconds", pins, seconds)
		return nil
	}
}
