package app_test

import (
	"context"
	"errors"
	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateProgram_Handle(t *testing.T) {
	ctx := context.Background()
	errTest := errors.New("")
	type args struct {
		ctx context.Context
		cmd cqs.Command
	}
	prog := fixtures.ProgramBuilder{}.Build()
	defaultArgs := args{
		ctx: ctx,
		cmd: app.CreateDailyProgramCommand{
			Program: &prog,
		},
	}
	tests := []struct {
		name string
		args args
		expectedErr, findErr,
		saveErr error
	}{
		{
			name: "with an invalid command, then it returns an invalid command error",
			args: args{
				ctx: ctx,
				cmd: invalidCommand{},
			},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and find hour returns a nil error, then it returns a create program error",
			args:        defaultArgs,
			expectedErr: app.CreateProgramError{},
		},
		{
			name:        "and find hour returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and save method returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     vo.NotFoundError{},
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "and save method returns nil, then it returns nil",
			args:    defaultArgs,
			findErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(`Given CreateDailyProgram command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &ProgramRepositoryMock{
				FindByHourFunc: func(ctx context.Context, hour *program.Hour) (*program.Program, error) {
					return nil, tt.findErr
				},
				SaveFunc: func(ctx context.Context, pr *program.Program) error {
					return tt.saveErr
				},
			}

			handler := app.NewCreateDailyProgram(repo)
			_, err := handler.Handle(tt.args.ctx, tt.args.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			}
		})
	}
}
