package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/stretchr/testify/require"
)

func TestCreateWeeklyProgram_Handle(t *testing.T) {
	errTest := errors.New("test error")
	ctx := context.Background()
	type args struct {
		ctx context.Context
		cmd cqs.Command
	}
	weekly := fixtures.WeeklyBuilder{}.Build()
	defaultArgs := args{
		ctx: ctx,
		cmd: app.CreateWeeklyProgramCommand{
			Weekly: &weekly,
		},
	}
	tests := []struct {
		name string
		args args
		findErr, saveErr,
		expectedErr error
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
			name:        "and find program returns a nil error, then it returns a create weekly program error",
			args:        defaultArgs,
			expectedErr: app.CreateWeeklyProgramError{},
		},
		{
			name:        "and find program returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and save program returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     vo.NotFoundError{},
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "and save program returns nil, then it returns nil",
			args:    defaultArgs,
			findErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a CreateWeeklyProgram command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &WeeklyProgramRepositoryMock{
				FindByDayFunc: func(ctx context.Context, day *program.WeekDay) (*program.Weekly, error) {
					return nil, tt.findErr
				},
				SaveFunc: func(ctx context.Context, programMoqParam *program.Weekly) error {
					return tt.saveErr
				},
			}
			handler := app.NewCreateWeeklyProgram(repo)
			_, err := handler.Handle(tt.args.ctx, tt.args.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			}
		})
	}
}
