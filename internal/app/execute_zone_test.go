package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/stretchr/testify/require"
)

func TestExecuteZoneHandle(t *testing.T) {
	errTest := errors.New("")
	zo := fixtures.ZoneBuilder{}.Build()
	tests := []struct {
		name                 string
		command              cqs.Command
		expectedErr, findErr error
		zone                 zone.Zone
	}{
		{
			name:        "and find zone returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and execute returns an error, then it returns an execute zone error",
			zone:        zo,
			expectedErr: app.ExecuteZoneError{},
			command: app.ExecuteZoneCmd{
				Seconds: 367,
				ZoneID:  "name",
			},
		},
		{
			name: "and execute returns an error, then it returns an execute zone error",
			zone: zo,
			command: app.ExecuteZoneCmd{
				Seconds: 36,
				ZoneID:  "name",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given an executeZone command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			zr := &ZoneRepositoryMock{
				FindByIDFunc: func(ctx context.Context, id string) (zone.Zone, error) {
					return tt.zone, tt.findErr
				},
			}
			handler := app.NewExecuteZone(zr)
			events, err := handler.Handle(context.Background(), tt.command)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Len(t, events, 1)
			ev := events[0]
			_, ok := ev.(zone.Executed)
			require.True(t, ok)
		})
	}
}
