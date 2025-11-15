package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestExecuteZoneWithStatusHandle(t *testing.T) {
	errTest := errors.New("")
	zo := fixtures.ZoneBuilder{}.Build()
	cmd := app.ExecuteZoneWithStatusCmd{
		Seconds: 367,
		ZoneID:  "name",
	}
	st := fixtures.StatusBuilder{Active: true}.Build()
	rainingWeather := fixtures.WeatherBuilder{Raining: true}.Build()
	statusRaining := fixtures.StatusBuilder{Weather: &rainingWeather}.Build()
	statusDeactivated := fixtures.StatusBuilder{Active: false}.Build()
	tests := []struct {
		name    string
		command cqs.Command
		expectedErr, findErr,
		stErr error
		zone          *zone.Zone
		expectedEvent string
		st            status.Status
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
			name:        "and status repository returns an error, then it returns same error",
			command:     cmd,
			zone:        zo,
			stErr:       errTest,
			expectedErr: errTest,
		},
		{
			name:        "and execute returns an error, then it returns an execute zone with status error",
			zone:        zo,
			st:          st,
			expectedErr: app.ExecuteZoneWithStatusError{},
			command:     cmd,
		},
		{
			name: "and execute nil, then it returns an executed event",
			zone: zo,
			st:   st,
			command: app.ExecuteZoneWithStatusCmd{
				Seconds: 36,
				ZoneID:  "name",
			},
			expectedEvent: zone.ExecutedEventName,
		},
		{
			name: "and its raining, then it returns an ignored event",
			zone: zo,
			st:   statusRaining,
			command: app.ExecuteZoneWithStatusCmd{
				Seconds: 36,
				ZoneID:  "name",
			},
			expectedEvent: zone.IgnoredEventName,
		},
		{
			name: "and system is deactivated, then it returns an ignored event",
			zone: zo,
			st:   statusDeactivated,
			command: app.ExecuteZoneWithStatusCmd{
				Seconds: 36,
				ZoneID:  "name",
			},
			expectedEvent: zone.IgnoredEventName,
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
			sr := &StatusRepositoryMock{}
			sr.FindFunc = func(ctx context.Context) (status.Status, error) {
				return tt.st, tt.stErr
			}
			handler := app.NewExecuteZoneWithStatus(zr, sr)
			events, err := handler.Handle(context.Background(), tt.command)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Len(t, events, 1)
			ev := events[0]
			require.Equal(t, tt.expectedEvent, ev.EventName())
		})
	}
}
