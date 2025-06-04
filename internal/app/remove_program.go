package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const (
	RemoveDailyProgramCommandName = "removeDailyProgram"
	RemoveOddProgramCommandName   = "removeOddProgram"
	RemoveEvenProgramCommandName  = "removeEvenProgram"
)

type RemoveDailyProgramCommand struct {
	Hour *program.Hour
}

func (c RemoveDailyProgramCommand) Name() string {
	return RemoveDailyProgramCommandName
}

type RemoveOddProgramCommand struct {
	Hour *program.Hour
}

func (c RemoveOddProgramCommand) Name() string {
	return RemoveOddProgramCommandName
}

type RemoveEvenProgramCommand struct {
	Hour *program.Hour
}

func (c RemoveEvenProgramCommand) Name() string {
	return RemoveEvenProgramCommandName
}

type RemoveDailyProgram struct {
	RemoveProgram
}

func (r RemoveDailyProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(RemoveDailyProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(RemoveDailyProgramCommandName, cmd.Name())
	}
	return nil, r.Remove(ctx, co.Hour)
}

func NewRemoveDailyProgram(pr ProgramRepository) *RemoveDailyProgram {
	return &RemoveDailyProgram{RemoveProgram{pr: pr}}
}

type RemoveOddProgram struct {
	RemoveProgram
}

func (r RemoveOddProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(RemoveOddProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(RemoveOddProgramCommandName, cmd.Name())
	}
	return nil, r.Remove(ctx, co.Hour)
}

func NewRemoveOddProgram(pr ProgramRepository) *RemoveOddProgram {
	return &RemoveOddProgram{RemoveProgram: RemoveProgram{
		pr: pr,
	}}
}

type RemoveEvenProgram struct {
	RemoveProgram
}

func (r RemoveEvenProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(RemoveEvenProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(RemoveEvenProgramCommandName, cmd.Name())
	}
	return nil, r.Remove(ctx, co.Hour)
}

func NewRemoveEvenProgram(pr ProgramRepository) *RemoveEvenProgram {
	return &RemoveEvenProgram{RemoveProgram: RemoveProgram{pr: pr}}
}

type RemoveProgram struct {
	pr ProgramRepository
}

func (r RemoveProgram) Remove(ctx context.Context, hour *program.Hour) error {
	if _, err := r.pr.FindByHour(ctx, hour); err != nil {
		return err
	}
	return r.pr.Remove(ctx, hour)
}
