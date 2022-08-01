package program_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestParseHour(t *testing.T) {
	tests := []struct {
		name, value string
	}{
		{
			name:  "with an invalid value, then it return an invalid execution hour error",
			value: "invalid",
		},
		{
			name:  "with a valid value, then it return a valid hour",
			value: "15:00",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Hour type,
		when ParseHour function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			hour, err := program.ParseHour(tt.value)
			if err != nil {
				require.ErrorIs(t, err, program.ErrInvalidExecutionHour)
				return
			}
			require.Equal(t, tt.value, hour.String())
		})
	}
}
