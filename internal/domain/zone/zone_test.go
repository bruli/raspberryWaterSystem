package zone_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
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
		tt := tt
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
