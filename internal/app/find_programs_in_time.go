package app

import (
	"context"
	"errors"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const FindProgramsInTimeQueryName = "findProgramsInTime"

type FindProgramsInTimeQuery struct {
	On          vo.Time
	Temperature float32
}

func (f FindProgramsInTimeQuery) Name() string {
	return FindProgramsInTimeQueryName
}

type FindProgramsInTime struct {
	Daily, Odd, Even ProgramRepository
	Weekly           WeeklyProgramRepository
	Temperature      TemperatureProgramRepository
}

func (f FindProgramsInTime) Handle(ctx context.Context, query cqs.Query) (any, error) {
	q, ok := query.(FindProgramsInTimeQuery)
	if !ok {
		return nil, cqs.NewInvalidQueryError(FindProgramsInTimeQueryName, query.Name())
	}
	hour, _ := program.ParseHour(q.On.HourStr())
	daily, err := f.findDaily(ctx, &hour)
	if err != nil {
		return nil, err
	}
	odd, err := f.findOdd(ctx, &hour)
	if err != nil {
		return nil, err
	}
	even, err := f.findEven(ctx, &hour)
	if err != nil {
		return nil, err
	}
	weekly, err := f.findWeekly(ctx, hour, time.Time(q.On).Weekday())
	if err != nil {
		return nil, err
	}
	temp, err := f.findTemperature(ctx, hour, q.Temperature)
	if err != nil {
		return nil, err
	}
	return ProgramsInTime{
		Daily:       daily,
		Odd:         odd,
		Even:        even,
		Weekly:      weekly,
		Temperature: temp,
	}, nil
}

func (f FindProgramsInTime) findDaily(ctx context.Context, hour *program.Hour) (*program.Program, error) {
	daily, err := f.Daily.FindByHour(ctx, hour)
	return f.findProgram(err, daily)
}

func (f FindProgramsInTime) findProgram(err error, prgInTime *program.Program) (*program.Program, error) {
	if err != nil {
		if !errors.As(err, &vo.NotFoundError{}) {
			return nil, err
		}
	}
	return prgInTime, nil
}

func (f FindProgramsInTime) findOdd(ctx context.Context, hour *program.Hour) (*program.Program, error) {
	odd, err := f.Odd.FindByHour(ctx, hour)
	return f.findProgram(err, odd)
}

func (f FindProgramsInTime) findEven(ctx context.Context, hour *program.Hour) (*program.Program, error) {
	even, err := f.Even.FindByHour(ctx, hour)
	return f.findProgram(err, even)
}

func (f FindProgramsInTime) findWeekly(ctx context.Context, hour program.Hour, day time.Weekday) (*program.Weekly, error) {
	weekDay := program.WeekDay(day)
	weekly, err := f.Weekly.FindByDayAndHour(ctx, &weekDay, &hour)
	if err != nil {
		if !errors.As(err, &vo.NotFoundError{}) {
			return nil, err
		}
	}
	return weekly, nil
}

func (f FindProgramsInTime) findTemperature(ctx context.Context, hour program.Hour, temperature float32) (*program.Temperature, error) {
	temp, err := f.Temperature.FindByTemperatureAndHour(ctx, temperature, hour)
	if err != nil {
		if !errors.As(err, &vo.NotFoundError{}) {
			return nil, err
		}
	}
	return &temp, nil
}

func NewFindProgramsInTime(
	daily, odd, even ProgramRepository,
	weekly WeeklyProgramRepository,
	temperature TemperatureProgramRepository,
) FindProgramsInTime {
	return FindProgramsInTime{Daily: daily, Odd: odd, Even: even, Weekly: weekly, Temperature: temperature}
}

type ProgramsInTime struct {
	Daily, Odd, Even *program.Program
	Weekly           *program.Weekly
	Temperature      *program.Temperature
}
