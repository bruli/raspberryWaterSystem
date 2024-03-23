package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestCreateProgramsHandle(t *testing.T) {
	errTest := errors.New("")
	cmd := app.CreateProgramsCmd{}
	tests := []struct {
		name string
		cmd  cqs.Command
		expectedErr, dailyErr,
		oddErr, evenErr,
		weeklyErr, tempErr error
	}{
		{
			name:        "with an invalid command, then it returns an invalid command error",
			cmd:         invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and daily repo returns error, then it returns same error",
			cmd:         cmd,
			dailyErr:    errTest,
			expectedErr: errTest,
		},
		{
			name:        "and odd repo returns error, then it returns same error",
			cmd:         cmd,
			oddErr:      errTest,
			expectedErr: errTest,
		},
		{
			name:        "and even repo returns error, then it returns same error",
			cmd:         cmd,
			evenErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and weekly repo returns error, then it returns same error",
			cmd:         cmd,
			weeklyErr:   errTest,
			expectedErr: errTest,
		},
		{
			name:        "and temperature repo returns error, then it returns same error",
			cmd:         cmd,
			tempErr:     errTest,
			expectedErr: errTest,
		},
		{
			name: "then it returns nil",
			cmd:  cmd,
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
			events, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Nil(t, events)
		})
	}
}
