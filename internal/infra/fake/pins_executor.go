package fake

import (
	"context"
	"time"
)

type PinsExecutor struct{}

func NewPinsExecutor() PinsExecutor {
	return PinsExecutor{}
}

func (p PinsExecutor) Execute(ctx context.Context, seconds uint, _ []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(time.Duration(seconds) * time.Second)
		return nil
	}
}
