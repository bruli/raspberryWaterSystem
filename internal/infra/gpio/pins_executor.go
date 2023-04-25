package gpio

import (
	"context"
	"fmt"
	"time"
)

type PinsExecutor struct {
	relays map[string]*pin
}

func NewPinsExecutor() *PinsExecutor {
	return &PinsExecutor{relays: relays}
}

func (p PinsExecutor) Execute(ctx context.Context, seconds uint, pins []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		activatedPins := make([]*pin, len(pins))
		for i, piNumber := range pins {
			activatePin, err := p.activatePin(piNumber)
			if err != nil {
				return err
			}
			activatedPins[i] = activatePin
		}
		time.Sleep(time.Duration(seconds) * time.Second)
		for _, act := range activatedPins {
			p.deActivatePin(act)
		}
		return nil
	}
}

func (p *PinsExecutor) activatePin(piNumber string) (*pin, error) {
	pi, ok := p.relays[piNumber]
	if !ok {
		return nil, InvalidPinToExecuteError{pinNumber: piNumber}
	}
	pi.output().low()
	return pi, nil
}

func (p PinsExecutor) deActivatePin(pi *pin) {
	pi.output().high()
}

type InvalidPinToExecuteError struct {
	pinNumber string
}

func (i InvalidPinToExecuteError) Error() string {
	return fmt.Sprintf("invalid pin %q to execute", i.pinNumber)
}
