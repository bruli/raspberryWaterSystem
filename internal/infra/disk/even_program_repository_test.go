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

func TestEvenProgramRepository(t *testing.T) {
	t.Run(`Given a EvenProgramRepository, `, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/even_programs.yml"
		populateFile(t, path)
		repo := disk.NewEvenProgramRepository(path)
		t.Run(`when Save method is called,
		then it save a odd program`, func(t *testing.T) {
			hour, err := program.ParseHour("12:15")
			require.NoError(t, err)
			hour2, err := program.ParseHour("10:30")
			require.NoError(t, err)
			odds := []program.Even{
				fixtures.EvenProgramBuilder{Hour: &hour, Zones: []string{
					"4",
				}}.Build(),
				fixtures.EvenProgramBuilder{Hour: &hour2, Zones: []string{
					"3",
				}}.Build(),
				fixtures.EvenProgramBuilder{Zones: []string{
					"1", "2",
				}}.Build(),
			}
			err = repo.Save(ctx, odds)
			require.NoError(t, err)
		})
		t.Run(`when FindAll method is called,
		then it returns a odd programs slice`, func(t *testing.T) {
			odds, err := repo.FindAll(ctx)
			require.NoError(t, err)
			require.Len(t, odds, 3)
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
			then it returns valid odd program`, func(t *testing.T) {
				hour, err := program.ParseHour("10:30")
				require.NoError(t, err)
				prg, err := repo.FindByHour(ctx, hour)
				require.NoError(t, err)
				require.Equal(t, hour, prg.Hour())
			})
		})
	})
}
