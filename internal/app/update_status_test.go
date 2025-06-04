package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/stretchr/testify/require"
)

func TestUpdateStatusHandle(t *testing.T) {
	errTest := errors.New("")
	currentSt := fixtures.StatusBuilder{}.Build()
	cmd := app.UpdateStatusCmd{Weather: fixtures.WeatherBuilder{}.Build()}
	tests := []struct {
		name string
		expectedErr, findErr,
		updateErr error
		st  status.Status
		cmd cqs.Command
	}{
		{
			name:        "with an invalid command, then it returns an invalid command error",
			cmd:         invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and find status returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
			cmd:         cmd,
		},
		{
			name:        "and update status returns an error, then it returns same error",
			st:          currentSt,
			updateErr:   errTest,
			expectedErr: errTest,
			cmd:         cmd,
		},
		{
			name: "and update status returns nil, then it returns nil",
			st:   currentSt,
			cmd:  cmd,
		},
	}
	for _, tt := range tests {
		t.Run(`Given UpdateStatus command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			sr := &StatusRepositoryMock{
				FindFunc: func(ctx context.Context) (status.Status, error) {
					return tt.st, tt.findErr
				},
				UpdateFunc: func(ctx context.Context, st status.Status) error {
					return tt.updateErr
				},
			}
			handler := app.NewUpdateStatus(sr)
			events, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Empty(t, events)
		})
	}
}
