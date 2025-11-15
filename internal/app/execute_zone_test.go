package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/internal/fixtures"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestExecuteZoneHandle(t *testing.T) {
	errTest := errors.New("")
	zo := fixtures.ZoneBuilder{}.Build()
	cmd := app.ExecuteZoneCmd{
		Seconds: 367,
		ZoneID:  "name",
	}
	tests := []struct {
		name                 string
		command              cqs.Command
		expectedErr, findErr error
		zone                 *zone.Zone
	}{
		{
			name:        "with invalid command, then it returns an invalid command error",
			command:     invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "and find zone returns an error, then it returns same error",
			command:     cmd,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and execute returns an error, then it returns an execute zone error",
			zone:        zo,
			expectedErr: app.ExecuteZoneError{},
			command:     cmd,
		},
		{
			name: "and execute nil, then it returns a valid event",
			zone: zo,
			command: app.ExecuteZoneCmd{
				Seconds: 36,
				ZoneID:  "name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(`Given an executeZone command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			zr := &ZoneRepositoryMock{
				FindByIDFunc: func(ctx context.Context, id string) (*zone.Zone, error) {
					return tt.zone, tt.findErr
				},
			}
			handler := app.NewExecuteZone(zr)
			events, err := handler.Handle(context.Background(), tt.command)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
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
