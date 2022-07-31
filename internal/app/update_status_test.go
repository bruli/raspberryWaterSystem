package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/stretchr/testify/require"
)

func TestUpdateStatusHandle(t *testing.T) {
	errTest := errors.New("")
	currentSt := fixtures.StatusBuilder{}.Build()
	tests := []struct {
		name string
		expectedErr, findErr,
		updateErr error
		st status.Status
	}{
		{
			name:        "and find status returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and update status returns an error, then it returns same error",
			st:          currentSt,
			updateErr:   errTest,
			expectedErr: errTest,
		},
		{
			name: "and update status returns nil, then it returns nil",
			st:   currentSt,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			events, err := handler.Handle(context.Background(), app.UpdateStatusCmd{Weather: fixtures.WeatherBuilder{}.Build()})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Empty(t, events)
		})
	}
}
