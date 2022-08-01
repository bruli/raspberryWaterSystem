package program_test

import (
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestNewDaily(t *testing.T) {
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
		when program function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			prg, err := program.NewDaily(tt.seconds, tt.hour, tt.zones)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.hour, prg.Hour())
			require.Equal(t, tt.seconds, prg.Seconds())
			require.Equal(t, tt.zones, prg.Zones())
		})
	}
}

func TestNewODD(t *testing.T) {
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
		when program function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			prg, err := program.NewOdd(tt.seconds, tt.hour, tt.zones)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.hour, prg.Hour())
			require.Equal(t, tt.seconds, prg.Seconds())
			require.Equal(t, tt.zones, prg.Zones())
		})
	}
}

func TestNewEven(t *testing.T) {
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
		when program function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			prg, err := program.NewEven(tt.seconds, tt.hour, tt.zones)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.hour, prg.Hour())
			require.Equal(t, tt.seconds, prg.Seconds())
			require.Equal(t, tt.zones, prg.Zones())
		})
	}
}

func TestNewWeekly(t *testing.T) {
	tests := []struct {
		name        string
		programs    []program.Program
		expectedErr error
	}{
		{
			name:        "with empty programs, then it returns an empty weekly programs error",
			expectedErr: program.ErrEmptyWeeklyPrograms,
		},
		{
			name: "with programs, then it returns an valid weekly",
			programs: []program.Program{
				fixtures.ProgramBuilder{}.Build(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Weekly struct,
		when NewWeekly function is called, `+tt.name, func(t *testing.T) {
			t.Parallel()
			week, err := program.NewWeekly(time.Friday, tt.programs)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, time.Friday, week.WeekDay())
			require.Equal(t, tt.programs, week.Programs())
		})
	}
}
