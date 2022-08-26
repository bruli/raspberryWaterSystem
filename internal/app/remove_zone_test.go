package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestRemoveZoneHandle(t *testing.T) {
	errTest := errors.New("")
	zo := fixtures.ZoneBuilder{}.Build()
	tests := []struct {
		name string
		expectedErr, findErr,
		removeErr error
		zone zone.Zone
	}{
		{
			name:        "and FindByID returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and remove method returns an error, then it returns same error",
			zone:        zo,
			removeErr:   errTest,
			expectedErr: errTest,
		},
		{
			name: "then it returns any error",
			zone: zo,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a RemoveZone command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			zr := &ZoneRepositoryMock{
				FindByIDFunc: func(ctx context.Context, id string) (zone.Zone, error) {
					return tt.zone, tt.findErr
				},
				RemoveFunc: func(ctx context.Context, zo zone.Zone) error {
					return tt.removeErr
				},
			}
			handler := app.NewRemoveZone(zr)
			events, err := handler.Handle(context.Background(), app.RemoveZoneCmd{})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Nil(t, events)
		})
	}
}
