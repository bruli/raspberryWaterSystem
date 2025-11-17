package status_test

import (
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/stretchr/testify/require"
)

func TestNewLight(t *testing.T) {
	sunrise := time.Date(2020, time.April, 14, 7, 45, 0, 0, time.UTC)
	sunset := time.Date(2020, time.April, 14, 17, 30, 0, 0, time.UTC)
	type args struct {
		sunrise time.Time
		sunset  time.Time
		now     time.Time
	}
	tests := []struct {
		name           string
		args           args
		expectedResult bool
		expectedError  error
	}{
		{
			name:          "with and invalid sunrise time, when it returns an invalid sunrise error",
			args:          args{},
			expectedError: status.ErrInvalidSunrise,
		},
		{
			name: "with and invalid sunset time, when it returns an invalid sunset error",
			args: args{
				sunrise: sunrise,
			},
			expectedError: status.ErrInvalidSunset,
		},
		{
			name: "with valid times and early hour, when it returns an is day false",
			args: args{
				sunrise: sunrise,
				sunset:  sunset,
				now:     time.Date(2020, time.April, 14, 5, 45, 0, 0, time.UTC),
			},
			expectedResult: false,
		},
		{
			name: "with valid times and later hour, when it returns an is day false",
			args: args{
				sunrise: sunrise,
				sunset:  sunset,
				now:     time.Date(2020, time.April, 14, 20, 45, 0, 0, time.UTC),
			},
			expectedResult: false,
		},
		{
			name: "with valid times and day hour, when it returns an is day true",
			args: args{
				sunrise: sunrise,
				sunset:  sunset,
				now:     time.Date(2020, time.April, 14, 10, 45, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a NewLight constructor,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := status.NewLight(tt.args.sunrise, tt.args.sunset)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.Equal(t, tt.expectedResult, got.IsDay(tt.args.now))
		})
	}
}
