package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/internal/app"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func TestFindAllProgramsHandle(t *testing.T) {
	errTest := errors.New("")
	dailies := []program.Program{
		fixtures.ProgramBuilder{}.Build(),
	}
	odds := []program.Program{
		fixtures.ProgramBuilder{}.Build(),
	}
	evens := []program.Program{
		fixtures.ProgramBuilder{}.Build(),
	}
	weeklies := []program.Weekly{
		fixtures.WeeklyBuilder{}.Build(),
	}
	temps := []program.Temperature{
		fixtures.TemperatureBuilder{}.Build(),
	}
	tests := []struct {
		name  string
		query cqs.Query
		expectedErr, dailyErr,
		oddErr, evenErr,
		weeklyErr, tempErr error
		expectedResult     any
		dailies, odd, even []program.Program
		weeklies           []program.Weekly
		temps              []program.Temperature
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
			dailies: []program.Program{
				fixtures.ProgramBuilder{}.Build(),
			},
		},
		{
			name:        "and even repository returns an error, then it returns same error",
			evenErr:     errTest,
			expectedErr: errTest,
			dailies:     dailies,
			odd:         odds,
		},
		{
			name:        "and weekly repository returns an error, then it returns same error",
			weeklyErr:   errTest,
			expectedErr: errTest,
			dailies:     dailies,
			odd:         odds,
			even:        evens,
		},
		{
			name:        "and temperature repository returns an error, then it returns same error",
			tempErr:     errTest,
			expectedErr: errTest,
			dailies:     dailies,
			odd:         odds,
			even:        evens,
			weeklies:    weeklies,
		},
		{
			name:     "then it returns a valid result",
			dailies:  dailies,
			odd:      odds,
			even:     evens,
			weeklies: weeklies,
			temps:    temps,
			expectedResult: app.AllPrograms{
				Daily:       dailies,
				Odd:         odds,
				Even:        evens,
				Weekly:      weeklies,
				Temperature: temps,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a FindAllPrograms query handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			temperature := &TemperatureProgramRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.Temperature, error) {
					return tt.temps, tt.tempErr
				},
			}
			daily := &ProgramRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.Program, error) {
					return tt.dailies, tt.dailyErr
				},
			}
			odd := &ProgramRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.Program, error) {
					return tt.odd, tt.oddErr
				},
			}
			even := &ProgramRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.Program, error) {
					return tt.even, tt.evenErr
				},
			}
			weekly := &WeeklyProgramRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.Weekly, error) {
					return tt.weeklies, tt.weeklyErr
				},
			}
			handler := app.NewFindAllPrograms(daily, odd, even, weekly, temperature)
			result, err := handler.Handle(context.Background(), tt.query)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
