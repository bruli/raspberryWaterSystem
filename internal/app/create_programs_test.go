package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestCreateProgramsHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name string
		expectedErr, dailyErr,
		oddErr, evenErr,
		weeklyErr, tempErr error
	}{
		{
			name:        "and daily repo returns error, then it returns same error",
			dailyErr:    errTest,
			expectedErr: errTest,
		},
		{
			name:        "and odd repo returns error, then it returns same error",
			oddErr:      errTest,
			expectedErr: errTest,
		},
		{
			name:        "and even repo returns error, then it returns same error",
			evenErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and weekly repo returns error, then it returns same error",
			weeklyErr:   errTest,
			expectedErr: errTest,
		},
		{
			name:        "and temperature repo returns error, then it returns same error",
			tempErr:     errTest,
			expectedErr: errTest,
		},
		{
			name: "then it returns nil",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a CreatePrograms command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			daily := &ProgramRepositoryMock{
				SaveFunc: func(ctx context.Context, programs []program.Program) error {
					return tt.dailyErr
				},
			}
			odd := &ProgramRepositoryMock{
				SaveFunc: func(ctx context.Context, programs []program.Program) error {
					return tt.oddErr
				},
			}
			even := &ProgramRepositoryMock{
				SaveFunc: func(ctx context.Context, programs []program.Program) error {
					return tt.evenErr
				},
			}
			weekly := &WeeklyProgramRepositoryMock{
				SaveFunc: func(ctx context.Context, programs []program.Weekly) error {
					return tt.weeklyErr
				},
			}
			temp := &TemperatureProgramRepositoryMock{
				SaveFunc: func(ctx context.Context, programs []program.Temperature) error {
					return tt.tempErr
				},
			}
			handler := app.NewCreatePrograms(daily, odd, even, weekly, temp)
			events, err := handler.Handle(context.Background(), app.CreateProgramsCmd{})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Nil(t, events)
		})
	}
}
