package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"

	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

func TestCreateStatusHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name                 string
		expectedErr, saveErr error
	}{
		{
			name:        "and save method returns an error, then it returns same error",
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name: "and save method returns nil, then it returns empty events",
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
