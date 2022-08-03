package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

const CreateProgramsCmdName = "createPrograms"

type CreateProgramsCmd struct {
	Daily, Odd, Even []program.Program
	Weekly           []program.Weekly
	Temperature      []program.Temperature
}

func (c CreateProgramsCmd) Name() string {
	return CreateProgramsCmdName
}

type CreatePrograms struct {
	daily, odd, even ProgramRepository
	weekly           WeeklyProgramRepository
	temperature      TemperatureProgramRepository
}

func (c CreatePrograms) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(CreateProgramsCmd)
	if err := c.daily.Save(ctx, co.Daily); err != nil {
		return nil, err
	}
	if err := c.odd.Save(ctx, co.Odd); err != nil {
		return nil, err
	}
	if err := c.even.Save(ctx, co.Even); err != nil {
		return nil, err
	}
	if err := c.weekly.Save(ctx, co.Weekly); err != nil {
		return nil, err
	}
	if err := c.temperature.Save(ctx, co.Temperature); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewCreatePrograms(
	daily, odd, even ProgramRepository,
	weekly WeeklyProgramRepository,
	temperature TemperatureProgramRepository,
) CreatePrograms {
	return CreatePrograms{
		daily:       daily,
		odd:         odd,
		even:        even,
		weekly:      weekly,
		temperature: temperature,
	}
}
