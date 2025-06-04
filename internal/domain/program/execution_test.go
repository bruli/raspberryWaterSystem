package program_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestNewExecution(t *testing.T) {
	tests := []struct {
		name  string
		zones []string
	}{
		{
			name: "with empty zones, then returns an empty execution zones error",
		},
		{
			name:  "with all values, then returns a valid struct",
			zones: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Execution struct,
		when NewExecution function is called `+tt.name, func(t *testing.T) {
			sec, err := program.ParseSeconds(20)
			require.NoError(t, err)
			exec, err := program.NewExecution(sec, tt.zones)
			if err != nil {
				require.ErrorIs(t, err, program.ErrEmptyExecutionZones)
				return
			}
			require.Equal(t, sec, exec.Seconds())
			require.Equal(t, tt.zones, exec.Zones())
		})
	}
}
