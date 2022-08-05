package fake

import (
	"context"
	"log"
	"time"
)

type PinsExecutor struct {
	log *log.Logger
}

func NewPinsExecutor(log *log.Logger) PinsExecutor {
	return PinsExecutor{log: log}
}

func (p PinsExecutor) Execute(ctx context.Context, seconds uint, pins []string) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	p.log.Printf("pins %s executed %v seconds", pins, seconds)
	return nil
}
