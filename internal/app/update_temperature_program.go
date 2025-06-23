package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const UpdateTemperatureProgramCommandName = "updateTemperatureProgram"

type UpdateTemperatureProgramCommand struct {
	Temperature float32
	Programs    []program.Program
}

func (u UpdateTemperatureProgramCommand) Name() string {
	return UpdateTemperatureProgramCommandName
}

type UpdateTemperatureProgram struct {
	repo TemperatureProgramRepository
}

func (u UpdateTemperatureProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(UpdateTemperatureProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(UpdateTemperatureProgramCommandName, cmd.Name())
	}
	pr, err := u.repo.FindByTemperature(ctx, co.Temperature)
	if err != nil {
		return nil, err
	}
	pr.Update(co.Programs)

	return nil, u.repo.Save(ctx, pr)
}

func NewUpdateTemperatureProgram(repo TemperatureProgramRepository) *UpdateTemperatureProgram {
	return &UpdateTemperatureProgram{repo: repo}
}
