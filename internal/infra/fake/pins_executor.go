package fake

import (
	"context"
	"log"
)

type PinsExecutor struct {
	log *log.Logger
}

func NewPinsExecutor(log *log.Logger) PinsExecutor {
	return PinsExecutor{log: log}
}

func (p PinsExecutor) Execute(ctx context.Context, seconds uint, pins []string) error {
	p.log.Printf("pins executed %v seconds", seconds)
	return nil
}
