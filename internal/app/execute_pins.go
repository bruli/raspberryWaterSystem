package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const ExecutePinsCmdName = "executePins"

type ExecutePinsCmd struct {
	Seconds uint
	Pins    []string
}

func (e ExecutePinsCmd) Name() string {
	return ExecutePinsCmdName
}

type ExecutePins struct {
	pe PinExecutor
}

func (e ExecutePins) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(ExecutePinsCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(ExecutePinsCmdName, cmd.Name())
	}
	return nil, e.pe.Execute(ctx, co.Seconds, co.Pins)
}

func NewExecutePins(pe PinExecutor) ExecutePins {
	return ExecutePins{pe: pe}
}
