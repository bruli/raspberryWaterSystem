package program_test

import (
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	hour, err := program.ParseHour("15:04")
	require.NoError(t, err)
	tests := []struct {
		name        string
		hour        program.Hour
		executions  []program.Execution
		expectedErr error
	}{
		{
			name:        "with empty executions, then it returns an empty programs error",
			hour:        hour,
			expectedErr: program.ErrEmptyPrograms,
		},
		{
			name: "with all values, then it returns a valid struct",
			hour: hour,
			executions: []program.Execution{
				fixtures.ExecutionBuilder{}.Build(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Program struct,
		when New function is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			prg, err := program.New(tt.hour, tt.executions)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.hour, prg.Hour())
			require.Equal(t, tt.executions, prg.Executions())
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
			expectedErr: program.ErrEmptyPrograms,
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
			week, err := program.NewWeekly(program.WeekDay(time.Friday), tt.programs)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, program.WeekDay(time.Friday), week.WeekDay())
			require.Equal(t, tt.programs, week.Programs())
		})
	}
}

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name     string
		programs []program.Program
	}{
		{
			name: "with any programs, then it returns an empty programs error",
		},
		{
			name: "with programs, then it returns a valid struct",
			programs: []program.Program{
				fixtures.ProgramBuilder{}.Build(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Temperature struct,
		when NewTemperature method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			tempValue := float32(20.5)
			temp, err := program.NewTemperature(tempValue, tt.programs)
			if err != nil {
				require.ErrorIs(t, err, program.ErrEmptyPrograms)
				return
			}
			require.Equal(t, tempValue, temp.Temperature())
			require.Equal(t, tt.programs, temp.Programs())
		})
	}
}
