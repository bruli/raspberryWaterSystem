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

func TestProgramRepository(t *testing.T) {
	t.Run(`Given a ProgramRepository, `, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/programs.yml"
		populateFile(t, path)
		repo := disk.NewProgramRepository(path)
		t.Run(`when Save method is called,
		then it save programs`, func(t *testing.T) {
			hour, err := program.ParseHour("12:15")
			require.NoError(t, err)
			hour2, err := program.ParseHour("10:30")
			require.NoError(t, err)
			exec := []program.Execution{
				fixtures.ExecutionBuilder{}.Build(),
				fixtures.ExecutionBuilder{}.Build(),
				fixtures.ExecutionBuilder{}.Build(),
			}
			programs := []program.Program{
				fixtures.ProgramBuilder{Hour: &hour, Executions: exec}.Build(),
				fixtures.ProgramBuilder{Hour: &hour2}.Build(),
				fixtures.ProgramBuilder{}.Build(),
			}
			err = repo.Save(ctx, programs)
			require.NoError(t, err)
		})
		t.Run(`when FindAll method is called,
		then it returns a programs slice`, func(t *testing.T) {
			dailies, err := repo.FindAll(ctx)
			require.NoError(t, err)
			require.Len(t, dailies, 3)
		})
		t.Run(`when FindByHour method is called `, func(t *testing.T) {
			t.Run(`with an nonexistent hour,
			then it returns a not found error`, func(t *testing.T) {
				hour, err := program.ParseHour("08:00")
				require.NoError(t, err)
				_, err = repo.FindByHour(ctx, hour)
				require.ErrorAs(t, err, &vo.NotFoundError{})
			})
			t.Run(`with an existent hour,
			then it returns valid program`, func(t *testing.T) {
				hour, err := program.ParseHour("12:15")
				require.NoError(t, err)
				prg, err := repo.FindByHour(ctx, hour)
				require.NoError(t, err)
				require.Equal(t, hour, prg.Hour())
				require.Len(t, prg.Executions(), 3)
			})
		})
	})
}
