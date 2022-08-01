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

func TestDailyProgramRepository(t *testing.T) {
	t.Run(`Given a DailyProgramRepository, `, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/daily_programs.yml"
		populateFile(t, path)
		repo := disk.NewDailyProgramRepository(path)
		t.Run(`when Save method is called,
		then it save a daily program`, func(t *testing.T) {
			hour, err := program.ParseHour("12:15")
			require.NoError(t, err)
			hour2, err := program.ParseHour("10:30")
			require.NoError(t, err)
			daily := []program.Daily{
				fixtures.DailyProgramBuilder{Hour: &hour, Zones: []string{
					"4",
				}}.Build(),
				fixtures.DailyProgramBuilder{Hour: &hour2, Zones: []string{
					"3",
				}}.Build(),
				fixtures.DailyProgramBuilder{Zones: []string{
					"1", "2",
				}}.Build(),
			}
			err = repo.Save(ctx, daily)
			require.NoError(t, err)
		})
		t.Run(`when FindAll method is called,
		then it returns a daily programs slice`, func(t *testing.T) {
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
			then it returns valid daily program`, func(t *testing.T) {
				hour, err := program.ParseHour("10:30")
				require.NoError(t, err)
				prg, err := repo.FindByHour(ctx, hour)
				require.NoError(t, err)
				require.Equal(t, hour, prg.Hour())
			})
		})
	})
}
