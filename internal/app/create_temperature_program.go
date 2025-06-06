package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const CreateTemperatureProgramCommandName = "createTemperatureProgram"

type CreateTemperatureProgramCommand struct {
	Temperature *program.Temperature
}

func (c CreateTemperatureProgramCommand) Name() string {
	return CreateTemperatureProgramCommandName
}

type CreateTemperatureProgram struct {
	repo TemperatureProgramRepository
}

func (c CreateTemperatureProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateTemperatureProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateTemperatureProgramCommandName, cmd.Name())
	}
	_, err := c.repo.FindByTemperature(ctx, co.Temperature.Temperature())
	switch {
	case err == nil:
		return nil, CreateTemperatureProgramError{msg: fmt.Sprintf("a temperature program with day %v already exists", co.Temperature)}
	case errors.As(err, &vo.NotFoundError{}):
		return nil, c.repo.Save(ctx, co.Temperature)
	default:
		return nil, err
	}
}

func NewCreateTemperatureProgram(repo TemperatureProgramRepository) *CreateTemperatureProgram {
	return &CreateTemperatureProgram{repo: repo}
}

type CreateTemperatureProgramError struct {
	msg string
}

func (c CreateTemperatureProgramError) Error() string {
	return fmt.Sprintf("failed to create a weekly program: %s", c.msg)
}
