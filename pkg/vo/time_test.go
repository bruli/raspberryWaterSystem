package vo_test

import (
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/stretchr/testify/require"
)

func TestParseFromTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		time     time.Time
		expected vo.Time
	}{
		{
			name: "with an empty time object, then it returns invalid zero time error",
			time: time.Time{},
		},
		{
			name:     "with a valid time object, then it returns a valid object",
			time:     now,
			expected: vo.Time(now),
		},
	}
	for _, tt := range tests {
		t.Run(`Given a parseFromTime function,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			voTime, err := vo.ParseFromTime(tt.time)
			if err != nil {
				require.ErrorIs(t, err, vo.ErrInvalidZeroTime)
				return
			}
			require.Equal(t, tt.expected, voTime)
			require.NotEqual(t, 0, len(voTime.EpochString()))
			require.NotEqual(t, 0, len(voTime.HourStr()))
			require.False(t, voTime.IsZero())
		})
	}
}

func TestParseFromEpochStr(t *testing.T) {
	tests := []struct {
		name, value string
	}{
		{
			name:  "with an invalid string, then it return an epoch string to time error",
			value: "invalid",
		},
		{
			name:  "with a valid string, then it return a valid time",
			value: "1234",
		},
	}
	for _, tt := range tests {
		t.Run(`Given a ParseFromEpochStr function,
		when is called `+tt.name, func(t *testing.T) {
			ti, err := vo.ParseFromEpochStr(tt.value)
			if err != nil {
				require.ErrorIs(t, err, vo.ErrEpochStrToTime)
				return
			}
			require.False(t, ti.IsZero())
		})
	}
}
