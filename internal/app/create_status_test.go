package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"

	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

func TestCreateStatusHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name string
		expectedErr, saveErr,
		findErr error
	}{
		{
			name:        "and find method returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and find method returns nil error, then it returns status already error",
			expectedErr: app.ErrStatusAlreadyExist,
		},
		{
			name:        "and save method returns an error, then it returns same error",
			findErr:     vo.NotFoundError{},
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "and save method returns nil, then it returns empty events",
			findErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		tt := tt
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
			handler := app.NewCreateStatus(sr)
			events, err := handler.Handle(context.Background(), app.CreateStatusCmd{
				StartedAt: time.Now(),
				Weather:   weather.New(20, 40, false),
			})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Nil(t, events)
		})
	}
}
