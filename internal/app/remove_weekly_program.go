package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const RemoveWeeklyProgramCommandName = "removeWeeklyProgram"

type RemoveWeeklyProgramCommand struct {
	Day *program.WeekDay
}

func (r RemoveWeeklyProgramCommand) Name() string {
	return RemoveWeeklyProgramCommandName
}

type RemoveWeeklyProgram struct {
	repo WeeklyProgramRepository
}

func (r RemoveWeeklyProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(RemoveWeeklyProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(RemoveWeeklyProgramCommandName, cmd.Name())
	}
	if _, err := r.repo.FindByDay(ctx, co.Day); err != nil {
		return nil, err
	}
	return nil, r.repo.Remove(ctx, co.Day)
}

func NewRemoveWeeklyProgram(repo WeeklyProgramRepository) *RemoveWeeklyProgram {
	return &RemoveWeeklyProgram{repo: repo}
}
