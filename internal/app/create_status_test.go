package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/internal/fixtures"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/stretchr/testify/require"
)

func TestCreateStatusHandle(t *testing.T) {
	errTest := errors.New("")
	cmd := app.CreateStatusCmd{
		StartedAt: vo.TimeNow(),
		Weather:   weather.New(20, 40, false),
	}
	light := fixtures.LightBuilder{}.Build()
	tests := []struct {
		name string
		expectedErr, saveErr,
		findErr, lightErr error
		light *status.Light
		cmd   cqs.Command
	}{
		{
			name:        "with an invalid command, then it returns an invalid command error",
			cmd:         invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and find method returns an error, then it returns same error",
			cmd:         cmd,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and find method returns nil error, then it returns status already error",
			cmd:         cmd,
			expectedErr: app.ErrStatusAlreadyExist,
		},
		{
			name:        "and find light method returns an error, then it returns same error",
			cmd:         cmd,
			findErr:     vo.NotFoundError{},
			lightErr:    errTest,
			expectedErr: errTest,
		},
		{
			name:        "and save method returns an error, then it returns same error",
			cmd:         cmd,
			findErr:     vo.NotFoundError{},
			light:       light,
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "and save method returns nil, then it returns empty events",
			cmd:     cmd,
			light:   light,
			findErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Create status command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			sr := &StatusRepositoryMock{
				SaveFunc: func(ctx context.Context, st status.Status) error {
					return tt.saveErr
				},
				FindFunc: func(ctx context.Context) (status.Status, error) {
					return status.Status{}, tt.findErr
				},
			}
			lr := &LightRepositoryMock{
				FindFunc: func(ctx context.Context, _ time.Time) (*status.Light, error) {
					return nil, tt.lightErr
				},
			}
			handler := app.NewCreateStatus(sr, lr)
			events, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Nil(t, events)
		})
	}
}
