package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

const FindAllProgramsQueryName = "findAllPrograms"

type FindAllProgramsQuery struct{}

func (f FindAllProgramsQuery) Name() string {
	return FindAllProgramsQueryName
}

type FindAllPrograms struct {
	daily, odd, even ProgramRepository
	weekly           WeeklyProgramRepository
	temperature      TemperatureProgramRepository
}

func (f FindAllPrograms) Handle(ctx context.Context, _ cqs.Query) (any, error) {
	dailies, err := f.daily.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	odd, err := f.odd.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	even, err := f.even.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	weekly, err := f.weekly.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	temp, err := f.temperature.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return AllPrograms{
		Daily:       dailies,
		Odd:         odd,
		Even:        even,
		Weekly:      weekly,
		Temperature: temp,
	}, nil
}

func NewFindAllPrograms(
	daily ProgramRepository,
	odd ProgramRepository,
	even ProgramRepository,
	weekly WeeklyProgramRepository,
	temperature TemperatureProgramRepository,
) FindAllPrograms {
	return FindAllPrograms{daily: daily, odd: odd, even: even, weekly: weekly, temperature: temperature}
}

type AllPrograms struct {
	Daily, Odd, Even []program.Program
	Weekly           []program.Weekly
	Temperature      []program.Temperature
}
