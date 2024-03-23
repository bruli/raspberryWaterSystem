package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestExecutePinsHandle(t *testing.T) {
	errTest := errors.New("")
	cmd := app.ExecutePinsCmd{}
	tests := []struct {
		name                 string
		execErr, expectedErr error
		cmd                  cqs.Command
	}{
		{
			name:        "with an invalid command, then it returns an invalid command error",
			cmd:         invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and executor returns an error, then it returns same error",
			cmd:         cmd,
			execErr:     errTest,
			expectedErr: errTest,
		},
		{
			name: "and executor returns nil, then it returns nil events",
			cmd:  cmd,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given an ExecutePins command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			pe := &PinExecutorMock{
				ExecuteFunc: func(ctx context.Context, seconds uint, pins []string) error {
					return tt.execErr
				},
			}
			handler := app.NewExecutePins(pe)
			events, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Nil(t, events)
		})
	}
}
