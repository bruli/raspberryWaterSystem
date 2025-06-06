package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const RemoveTemperatureProgramCommandName = "removeTemperatureProgram"

type RemoveTemperatureProgramCommand struct {
	Temperature float32
}

func (r RemoveTemperatureProgramCommand) Name() string {
	return RemoveTemperatureProgramCommandName
}

type RemoveTemperatureProgram struct {
	repo TemperatureProgramRepository
}

func (r RemoveTemperatureProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(RemoveTemperatureProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(RemoveTemperatureProgramCommandName, cmd.Name())
	}
	if _, err := r.repo.FindByTemperature(ctx, co.Temperature); err != nil {
		return nil, err
	}
	return nil, r.repo.Remove(ctx, co.Temperature)
}

func NewRemoveTemperatureProgram(repo TemperatureProgramRepository) *RemoveTemperatureProgram {
	return &RemoveTemperatureProgram{repo: repo}
}
