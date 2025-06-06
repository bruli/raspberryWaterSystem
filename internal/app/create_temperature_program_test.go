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

func TestCreateTemperatureProgram_Handle(t *testing.T) {
	errTest := errors.New("test error")
	ctx := context.Background()
	type args struct {
		ctx context.Context
		cmd cqs.Command
	}
	weekly := fixtures.TemperatureBuilder{}.Build()
	defaultArgs := args{
		ctx: ctx,
		cmd: app.CreateTemperatureProgramCommand{
			Temperature: &weekly,
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
			name:        "and find program returns a nil error, then it returns a create remove program error",
			args:        defaultArgs,
			expectedErr: app.CreateTemperatureProgramError{},
		},
		{
			name:        "and find program returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and remove program returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     vo.NotFoundError{},
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "and remove program returns nil, then it returns nil",
			args:    defaultArgs,
			findErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a CreateTemperatureProgram command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &TemperatureProgramRepositoryMock{
				FindByTemperatureFunc: func(ctx context.Context, temperature float32) (*program.Temperature, error) {
					return nil, tt.findErr
				},
				SaveFunc: func(ctx context.Context, program *program.Temperature) error {
					return tt.saveErr
				},
			}
			handler := app.NewCreateTemperatureProgram(repo)
			_, err := handler.Handle(tt.args.ctx, tt.args.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			}
		})
	}
}
