package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const CreateProgramCommandName = "createProgram"

type CreateProgramCommand struct {
	Program *program.Program
}

func (c CreateProgramCommand) Name() string {
	return CreateProgramCommandName
}

type CreateProgram struct {
	repo ProgramRepository
}

func (c CreateProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateProgramCommandName, cmd.Name())
	}
	hour := co.Program.Hour()
	_, err := c.repo.FindByHour(ctx, &hour)
	switch {
	case err == nil:
		return nil, CreateProgramError{msg: fmt.Sprintf("a program with hour %s, already exists", hour.String())}
	case errors.As(err, &vo.NotFoundError{}):
		return nil, c.repo.Save(ctx, co.Program)
	default:
		return nil, err
	}
}

func NewCreateProgram(repo ProgramRepository) *CreateProgram {
	return &CreateProgram{repo: repo}
}

type CreateProgramError struct {
	msg string
}

func (c CreateProgramError) Error() string {
	return fmt.Sprintf("failed to create program: %s", c.msg)
}
