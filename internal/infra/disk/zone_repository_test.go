//+build infra

package disk_test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/google/uuid"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/stretchr/testify/require"
)

func TestZoneRepository(t *testing.T) {
	t.Run(`Given a ZoneRepository, `, func(t *testing.T) {
		ctx := context.Background()
		repo := disk.NewZoneRepository("/tmp/zones.yml")
		var savedZone *zone.Zone
		_ = savedZone
		t.Run(`when Save method is called,
		then it save a zone`, func(t *testing.T) {
			zo := fixtures.ZoneBuilder{}.Build()
			err := repo.Save(ctx, zo)
			require.NoError(t, err)
			savedZone = &zo
		})
		t.Run(`when Save method is called
		with a duplicated zone,
		then it return a duplicated error`, func(t *testing.T) {
			err := repo.Save(ctx, *savedZone)
			require.ErrorAs(t, err, &disk.DuplicatedZoneError{})
		})
		t.Run(`when FindByID method is called
		with a valid id,
		then it return a zone`, func(t *testing.T) {
			id := savedZone.Id()
			zo, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, *savedZone, zo)
		})
		t.Run(`when FindByID method is called
		with an invalid id,
		then it return a not found error`, func(t *testing.T) {
			_, err := repo.FindByID(ctx, uuid.New().String())
			require.ErrorAs(t, err, &vo.NotFoundError{})
		})
		t.Run(`when update method is called
		with an invalid id,
		then it return a not found error`, func(t *testing.T) {
			zo := fixtures.ZoneBuilder{}.Build()
			err := repo.Update(ctx, zo)
			require.ErrorAs(t, err, &vo.NotFoundError{})
		})
		t.Run(`when update method is called,
		then it update zone`, func(t *testing.T) {
			updated := *savedZone
			updated.Update(spew.Sprintf("pepito%s", uuid.New().String()), savedZone.Relays())
			err := repo.Update(ctx, updated)
			require.NoError(t, err)

			updatedZone, err := repo.FindByID(ctx, updated.Id())
			require.NoError(t, err)
			require.Equal(t, updated, updatedZone)
		})
	})
}
