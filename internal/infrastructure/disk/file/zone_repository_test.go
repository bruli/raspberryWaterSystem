package file

import (
	"errors"
	zone2 "github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestZoneRepository_GetZones(t *testing.T) {
	validZon, err := zone2.New("bb", "bb", []string{"1"})
	assert.NoError(t, err)
	validZons := zone2.Zones{}
	validZons.Add(*validZon)

	zons := zones{}
	zon := newZone("bb", "bb", []string{"1"})
	zons.add(zon)
	validData, err := yaml.Marshal(&zons)
	assert.NoError(t, err)

	tests := map[string]struct {
		zo   *zone2.Zones
		data []byte
		err  error
	}{
		"it should return empty zones when reader returns error": {
			zo:  &zone2.Zones{},
			err: errors.New("error"),
		},
		"it should return empty zones when unmarshal returns error": {
			zo:   &zone2.Zones{},
			data: []byte("invalid"),
		},
		"it should return zones": {
			zo:   &validZons,
			data: validData,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			read := readerMock{}
			zoneRepo := ZoneRepository{repository: &repository{reader: &read}}

			read.readFunc = func() ([]byte, error) {
				return tt.data, tt.err
			}
			zons := zoneRepo.GetZones()
			assert.Equal(t, tt.zo, zons)
		})
	}
}
