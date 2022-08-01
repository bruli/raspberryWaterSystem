package program_test

import (
	"testing"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	seconds, err := program.ParseSeconds(20)
	require.NoError(t, err)
	zones := []string{"1", "2"}
	hour, err := program.ParseHour("15:04")
	require.NoError(t, err)
	tests := []struct {
		name        string
		seconds     program.Seconds
		hour        program.Hour
		zones       []string
		expectedErr error
	}{
		{
			name:        "with an invalid seconds, then it returns an zero program seconds error",
			seconds:     program.Seconds(time.Duration(0) * time.Second),
			zones:       zones,
			hour:        hour,
			expectedErr: program.ErrZeroProgramSeconds,
		},
		{
			name:        "with an empty zones, then it returns an empty execution zones error",
			seconds:     seconds,
			hour:        hour,
			expectedErr: program.ErrEmptyExecutionZones,
		},
		{
			name:    "with all values, then it returns a program",
			seconds: seconds,
			hour:    hour,
			zones:   zones,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Program struct,
		when New function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			prg, err := program.New(tt.seconds, tt.hour, tt.zones)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.hour, prg.Execution().Hour())
			require.Equal(t, tt.seconds, prg.Seconds())
			require.Equal(t, tt.zones, prg.Execution().Zones())
		})
	}
}
