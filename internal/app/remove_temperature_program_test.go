package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestRemoveTemperature_Handle(t *testing.T) {
	errTest := errors.New("")
	type args struct {
		ctx context.Context
		cmd cqs.Command
	}
	defaultArgs := args{
		ctx: t.Context(),
		cmd: app.RemoveTemperatureProgramCommand{},
	}
	tests := []struct {
		name string
		args args
		findErr, removeErr,
		expectedErr error
	}{
		{
			name: "with an invalid command, then it returns an invalid command error",
			args: args{
				ctx: t.Context(),
				cmd: invalidCommand{},
			},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and find program returns an error, then it returns same error",
			args:        defaultArgs,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and Temperature remove program returns an error, then it returns same error",
			args:        defaultArgs,
			removeErr:   errTest,
			expectedErr: errTest,
		},
		{
			name: "and remove program returns nil, then it returns nil",
			args: defaultArgs,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a RemoveTemperatureProgram command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &TemperatureProgramRepositoryMock{
				FindByTemperatureFunc: func(ctx context.Context, temperature float32) (*program.Temperature, error) {
					return nil, tt.findErr
				},
				RemoveFunc: func(ctx context.Context, temperature float32) error {
					return tt.removeErr
				},
			}
			handler := app.NewRemoveTemperatureProgram(repo)
			_, err := handler.Handle(tt.args.ctx, tt.args.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			}
		})
	}
}
