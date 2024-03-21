package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestFindProgramsInTimeHandle(t *testing.T) {
	errTest := errors.New("")
	daily := fixtures.ProgramBuilder{}.Build()
	odd := fixtures.ProgramBuilder{}.Build()
	even := fixtures.ProgramBuilder{}.Build()
	weekly := fixtures.WeeklyBuilder{}.Build()
	temp := fixtures.TemperatureBuilder{}.Build()
	tests := []struct {
		name string
		expectedErr, dailyErr,
		oddErr, evenErr,
		weeklyErr, tempErr error
		expectedResult   any
		daily, odd, even program.Program
		weekly           program.Weekly
		temp             program.Temperature
	}{
		{
			name:        "and daily repository returns an error, then it returns same error",
			dailyErr:    errTest,
			expectedErr: errTest,
		},
		{
			name:        "and odd repository returns an error, then it returns same error",
			oddErr:      errTest,
			expectedErr: errTest,
			daily:       daily,
		},
		{
			name:        "and even repository returns an error, then it returns same error",
			evenErr:     errTest,
			expectedErr: errTest,
			daily:       daily,
			odd:         odd,
		},
		{
			name:        "and weekly repository returns an error, then it returns same error",
			weeklyErr:   errTest,
			expectedErr: errTest,
			daily:       daily,
			odd:         odd,
			even:        even,
		},
		{
			name:        "and temperature repository returns an error, then it returns same error",
			tempErr:     errTest,
			expectedErr: errTest,
			daily:       daily,
			odd:         odd,
			even:        even,
			weekly:      weekly,
		},
		{
			name:   "then it returns a valid result",
			daily:  daily,
			odd:    odd,
			even:   even,
			weekly: weekly,
			temp:   temp,
			expectedResult: app.ProgramsInTime{
				Daily:       &daily,
				Odd:         &odd,
				Even:        &even,
				Weekly:      &weekly,
				Temperature: &temp,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a FindProgramsInTime query handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			daily := &ProgramRepositoryMock{
				FindByHourFunc: func(ctx context.Context, hour program.Hour) (program.Program, error) {
					return tt.daily, tt.dailyErr
				},
			}
			odd := &ProgramRepositoryMock{
				FindByHourFunc: func(ctx context.Context, hour program.Hour) (program.Program, error) {
					return tt.odd, tt.oddErr
				},
			}
			even := &ProgramRepositoryMock{
				FindByHourFunc: func(ctx context.Context, hour program.Hour) (program.Program, error) {
					return tt.even, tt.evenErr
				},
			}
			weekly := &WeeklyProgramRepositoryMock{
				FindByDayAndHourFunc: func(ctx context.Context, day program.WeekDay, hour program.Hour) (program.Weekly, error) {
					return tt.weekly, tt.weeklyErr
				},
			}
			temperature := &TemperatureProgramRepositoryMock{
				FindByTemperatureAndHourFunc: func(ctx context.Context, temperature float32, hour program.Hour) (program.Temperature, error) {
					return tt.temp, tt.tempErr
				},
			}
			handler := app.NewFindProgramsInTime(daily, odd, even, weekly, temperature)
			result, err := handler.Handle(context.Background(), app.FindProgramsInTimeQuery{})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
