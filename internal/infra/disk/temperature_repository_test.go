//go:build infra

package disk_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/stretchr/testify/require"
)

func TestTemperatureRepository(t *testing.T) {
	t.Run(`Given a Temperature repository,`, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/temperature_programs.yml"
		defer populateFile(t, path)
		repo := disk.NewTemperatureProgramRepository(path)
		t.Run(`when Save method is called,
		then it save temperature programs`, func(t *testing.T) {
			hour1, err := program.ParseHour("20:12")
			require.NoError(t, err)
			hour2, err := program.ParseHour("21:12")
			require.NoError(t, err)
			hour3, err := program.ParseHour("22:12")
			require.NoError(t, err)
			seconds1, err := program.ParseSeconds(15)
			require.NoError(t, err)
			executions := []program.Execution{
				fixtures.ExecutionBuilder{Zones: []string{"1"}, Seconds: &seconds1}.Build(),
				fixtures.ExecutionBuilder{Zones: []string{"2"}}.Build(),
			}
			programs1 := []program.Program{
				fixtures.ProgramBuilder{Hour: &hour1}.Build(),
				fixtures.ProgramBuilder{Executions: executions}.Build(),
			}
			programs2 := []program.Program{
				fixtures.ProgramBuilder{Hour: &hour2, Executions: []program.Execution{
					fixtures.ExecutionBuilder{}.Build(),
				}}.Build(),
			}
			programs3 := []program.Program{
				fixtures.ProgramBuilder{Hour: &hour3, Executions: []program.Execution{
					fixtures.ExecutionBuilder{}.Build(),
					fixtures.ExecutionBuilder{}.Build(),
					fixtures.ExecutionBuilder{}.Build(),
				}}.Build(),
			}
			temperatures := []program.Temperature{
				fixtures.TemperatureBuilder{Programs: programs1, Temperature: vo.Float32Ptr(20.3)}.Build(),
				fixtures.TemperatureBuilder{Programs: programs2, Temperature: vo.Float32Ptr(19.3)}.Build(),
				fixtures.TemperatureBuilder{Programs: programs3, Temperature: vo.Float32Ptr(22.3)}.Build(),
			}
			for _, temp := range temperatures {
				err = repo.Save(ctx, &temp)
				require.NoError(t, err)

			}
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
			t.Run(`with an major temperature and hour,
			then it returns a valid temperature program`, func(t *testing.T) {
				hour, err := program.ParseHour("15:10")
				require.NoError(t, err)
				tempValue := float32(21)
				temp, err := repo.FindByTemperatureAndHour(ctx, tempValue, hour)
				require.NoError(t, err)
				require.Len(t, temp.Programs(), 2)
				require.Len(t, temp.Programs()[0].Executions(), 1)
				require.Equal(t, []string{"1"}, temp.Programs()[0].Executions()[0].Zones())
				require.Equal(t, 15, temp.Programs()[0].Executions()[0].Seconds().Int())
				require.Len(t, temp.Programs()[1].Executions(), 1)
				require.Equal(t, []string{"2"}, temp.Programs()[1].Executions()[0].Zones())
				require.Equal(t, 20, temp.Programs()[1].Executions()[0].Seconds().Int())
				require.Equal(t, tempValue, temp.Temperature())
			})
		})
		t.Run(`when FindByTemperature method is called,`, func(t *testing.T) {
			t.Run(`with an invalid temperature,
			then it returns a not found error`, func(t *testing.T) {
				_, err := repo.FindByTemperature(ctx, float32(40))
				require.ErrorAs(t, err, &vo.NotFoundError{})
			})
			t.Run(`with a valid temperature,
			then it returns a valid temperature program`, func(t *testing.T) {
				tempValue := float32(20.3)
				find, err := repo.FindByTemperature(ctx, tempValue)
				require.NoError(t, err)
				require.NotNil(t, find)
			})
		})
	})
}
