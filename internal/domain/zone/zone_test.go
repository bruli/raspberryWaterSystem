package zone_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/internal/fixtures"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	rel1, err := zone.ParseRelay(zone.OneRelayID)
	require.NoError(t, err)
	relays := []zone.Relay{
		rel1,
	}
	tests := []struct {
		name, id, zoneName string
		relays             []zone.Relay
		expectedErr        error
	}{
		{
			name:        "with an empty id, then it return an invalid zone id",
			expectedErr: zone.ErrInvalidZoneID,
		},
		{
			name:        "with an empty name, then it return an invalid zone name",
			id:          "id",
			expectedErr: zone.ErrInvalidZoneName,
		},
		{
			name:        "with an empty relays, then it return an invalid zone relays",
			id:          "id",
			zoneName:    "name",
			expectedErr: zone.ErrInvalidZoneRelays,
		},
		{
			name:     "with all values, then it return a valid zone",
			id:       "id",
			zoneName: "name",
			relays:   relays,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Zone struct,
		when the constructor is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			zo, err := zone.New(tt.id, tt.zoneName, tt.relays)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.Equal(t, zo.Id(), tt.id)
			require.Equal(t, zo.Name(), tt.zoneName)
			require.Equal(t, zo.Relays(), tt.relays)
		})
	}
}

func TestZoneExecute(t *testing.T) {
	relays := []zone.Relay{
		fixtures.RelayBuilder(1),
		fixtures.RelayBuilder(2),
	}
	tests := []struct {
		name    string
		relay   []zone.Relay
		seconds uint
	}{
		{
			name:    "with invalid seconds, then it returns an invalid seconds execution zone error",
			seconds: 505,
		},
		{
			name:    "with valid seconds, then it returns a valid event",
			seconds: 200,
			relay:   relays,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a built zone struct,
		when Execute method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			var zo zone.Zone
			zo.Hydrate("bbf", "Bonsai", tt.relay)
			err := zo.Execute(tt.seconds)
			if err != nil {
				require.ErrorIs(t, err, zone.ErrInvalidSecondsExecutionZone)
				return
			}
			ev := zo.Events()[0]
			execEv, _ := ev.(zone.Executed)
			require.Equal(t, []string{"18", "17"}, execEv.RelayPins)
			require.Equal(t, uint(200), execEv.Seconds)
		})
	}
}
