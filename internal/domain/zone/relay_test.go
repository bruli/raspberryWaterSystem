package zone_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/stretchr/testify/require"
)

func TestParseRelay(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		expectedRelay zone.RelayID
	}{
		{
			name: "with an unknown id, then it return a unknown relay error",
			id:   99,
		},
		{
			name:          "with an unknown id, then it return a unknown valid relay id",
			id:            1,
			expectedRelay: zone.OneRelayID,
		},
	}
	for _, tt := range tests {

		t.Run(`Given a ParseRelay function,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			relay, err := zone.ParseRelay(tt.id)
			if err != nil {
				require.ErrorIs(t, err, zone.ErrUnknownRelay)
				return
			}
			require.Equal(t, tt.expectedRelay, relay.Id())
		})
	}
}
