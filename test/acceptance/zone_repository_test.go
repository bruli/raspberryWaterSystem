package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/disk/file"
	zone2 "github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZoneRepository(t *testing.T) {
	config := getConfig()
	z, err := zone2.NewZonesStub()
	assert.NoError(t, err)
	repo := file.NewZoneRepository(config.ZonesFile)
	err = repo.Save(z)

	assert.NoError(t, err)

	readZones := repo.GetZones()
	assert.Equal(t, 1, len(*readZones))

	zone := z[0]
	foundedZone := repo.Find(zone.Id())
	assert.Equal(t, &zone, foundedZone)
}
