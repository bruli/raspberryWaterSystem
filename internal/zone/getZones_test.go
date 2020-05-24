package zone_test

import (
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetZones_Get(t *testing.T) {
	zo := zone.Zones{}
	z, _ := zone.New("1", "name", []string{"1", "2", "3"})
	zo.Add(*z)
	zoneRepositoryM := zone.RepositoryMock{}
	getZ := zone.NewGetter(&zoneRepositoryM)
	t.Run("it should return zones", func(t *testing.T) {
		zoneRepositoryM.GetZonesFunc = func() *zone.Zones {
			return &zo
		}

		r := getZ.Get()
		assert.Equal(t, len(zo), len(*r))
	})
}
