package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const CreateDailyProgramCommandName = "createDailyProgram"

type CreateDailyProgramCommand struct {
	Program *program.Program
}

func (c CreateDailyProgramCommand) Name() string {
	return CreateDailyProgramCommandName
}

type CreateDailyProgram struct {
	CreateProgram
}

func (c CreateDailyProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateDailyProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateDailyProgramCommandName, cmd.Name())
	}
	return nil, c.Create(ctx, co.Program)
}

func NewCreateDailyProgram(repo ProgramRepository) *CreateDailyProgram {
	return &CreateDailyProgram{
		CreateProgram: CreateProgram{repo: repo},
	}
}

type CreateProgramError struct {
	msg string
}

func (c CreateProgramError) Error() string {
	return fmt.Sprintf("failed to create program: %s", c.msg)
}

type CreateProgram struct {
	repo ProgramRepository
}

func (p CreateProgram) Create(ctx context.Context, program *program.Program) error {
	hour := program.Hour()
	_, err := p.repo.FindByHour(ctx, &hour)
	switch {
	case err == nil:
		return CreateProgramError{msg: fmt.Sprintf("a program with hour %s, already exists", hour.String())}
	case errors.As(err, &vo.NotFoundError{}):
		return p.repo.Save(ctx, program)
	default:
		return err
	}
}
