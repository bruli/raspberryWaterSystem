package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const CreateWeeklyProgramCommandName = "createWeeklyProgram"

type CreateWeeklyProgramCommand struct {
	Weekly *program.Weekly
}

func (c CreateWeeklyProgramCommand) Name() string {
	return CreateWeeklyProgramCommandName
}

type CreateWeeklyProgram struct {
	repo WeeklyProgramRepository
}

func (c CreateWeeklyProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateWeeklyProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateWeeklyProgramCommandName, cmd.Name())
	}
	day := co.Weekly.WeekDay()
	_, err := c.repo.FindByDay(ctx, &day)
	switch {
	case err == nil:
		return nil, CreateWeeklyProgramError{msg: fmt.Sprintf("a weekly program with day %s already exists", day.String())}
	case errors.As(err, &vo.NotFoundError{}):
		return nil, c.repo.Save(ctx, co.Weekly)
	default:
		return nil, err
	}
}

func NewCreateWeeklyProgram(repo WeeklyProgramRepository) *CreateWeeklyProgram {
	return &CreateWeeklyProgram{repo: repo}

}

type CreateWeeklyProgramError struct {
	msg string
}

func (c CreateWeeklyProgramError) Error() string {
	return fmt.Sprintf("failed to create a weekly program: %s", c.msg)
}
