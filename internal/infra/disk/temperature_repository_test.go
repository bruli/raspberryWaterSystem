//go:build infra
// +build infra

package disk_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/stretchr/testify/require"
)

func TestTemperatureRepository(t *testing.T) {
	t.Run(`Given a Temperature repository,`, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/temperature_programs.yml"
		populateFile(t, path)
		repo := disk.NewTemperatureProgramRepository(path)
		t.Run(`when Save method is called,
		then it save temperature programs`, func(t *testing.T) {
			hour, err := program.ParseHour("20:12")
			require.NoError(t, err)
			executions := []program.Execution{
				fixtures.ExecutionBuilder{}.Build(),
				fixtures.ExecutionBuilder{}.Build(),
				fixtures.ExecutionBuilder{}.Build(),
			}
			programs := []program.Program{
				fixtures.ProgramBuilder{Hour: &hour}.Build(),
				fixtures.ProgramBuilder{Executions: executions}.Build(),
			}
			temperatures := []program.Temperature{
				fixtures.TemperatureBuilder{Programs: programs, Temperature: vo.Float32Ptr(20.3)}.Build(),
				fixtures.TemperatureBuilder{Programs: programs, Temperature: vo.Float32Ptr(19.3)}.Build(),
				fixtures.TemperatureBuilder{Programs: programs, Temperature: vo.Float32Ptr(22.3)}.Build(),
			}
			err = repo.Save(ctx, temperatures)
			require.NoError(t, err)
		})
		t.Run(`when FindAll method is called,
		then it returns a temperature programs slice`, func(t *testing.T) {
			prgms, err := repo.FindAll(ctx)
			require.NoError(t, err)
			require.Len(t, prgms, 3)
		})
		t.Run(`when FindByTemperatureAndHour method is called,`, func(t *testing.T) {
			t.Run(`with an invalid temperature,
			then it returns a not found error`, func(t *testing.T) {
				hour, err := program.ParseHour("08:00")
				require.NoError(t, err)
				_, err = repo.FindByTemperatureAndHour(ctx, float32(40), hour)
				require.ErrorAs(t, err, &vo.NotFoundError{})
			})
			t.Run(`with an invalid hour,
			then it returns a not found error`, func(t *testing.T) {
				hour, err := program.ParseHour("08:00")
				require.NoError(t, err)
				_, err = repo.FindByTemperatureAndHour(ctx, float32(22.3), hour)
				require.ErrorAs(t, err, &vo.NotFoundError{})
			})
			t.Run(`with an valid temperature and hour,
			then it returns a valid temperature program`, func(t *testing.T) {
				hour, err := program.ParseHour("15:10")
				require.NoError(t, err)
				tempValue := float32(22.3)
				temp, err := repo.FindByTemperatureAndHour(ctx, tempValue, hour)
				require.NoError(t, err)
				require.Len(t, temp.Programs(), 3)
				require.Equal(t, tempValue, temp.Temperature())
			})
		})
	})
}
