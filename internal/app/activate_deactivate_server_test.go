package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestActivateDeactivateServerHandle(t *testing.T) {
	errTest := errors.New("")
	st := fixtures.StatusBuilder{}.Build()
	tests := []struct {
		name string
		cmd  cqs.Command
		expectedErr, findErr,
		updateErr error
		status status.Status
	}{
		{
			name:        "with an invalid command, then it returns an invalid command error",
			cmd:         &invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "with a valid command and find method returns an error, then it returns same error",
			cmd:         app.ActivateDeactivateServerCmd{Active: true},
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "with a valid command and update method returns an error, then it returns same error",
			cmd:         app.ActivateDeactivateServerCmd{Active: true},
			status:      st,
			updateErr:   errTest,
			expectedErr: errTest,
		},
		{
			name:   "with a valid command and update method returns nil, then it update status",
			cmd:    app.ActivateDeactivateServerCmd{Active: false},
			status: st,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given ActivateDeactivateServer command handler,
		when handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			stRepo := &StatusRepositoryMock{}
			stRepo.FindFunc = func(ctx context.Context) (status.Status, error) {
				return tt.status, tt.findErr
			}
			stRepo.UpdateFunc = func(ctx context.Context, st status.Status) error {
				return tt.updateErr
			}
			handler := app.NewActivateDeactivateServer(stRepo)
			_, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

type invalidCommand struct{}

func (i invalidCommand) Name() string {
	return "invalid"
}
