package zone_test

import (
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemover_Remove(t *testing.T) {
	zon, err := zone.New("1", "test", []string{"1"})
	zons := zone.Zones{}
	zons.Add(*zon)
	assert.NoError(t, err)
	tests := map[string]struct {
		zoneID           string
		zon              *zone.Zone
		zons             *zone.Zones
		err, expectedErr error
	}{
		"it should return error when zone does not exists": {
			zoneID:      "1",
			expectedErr: zone.NewNotFound("1"),
		},
		"it should return error when save returns error": {
			zoneID:      "1",
			zon:         zon,
			zons:        &zons,
			err:         errors.New("error"),
			expectedErr: fmt.Errorf("failed to remove zoneID %s: %s", "1", errors.New("error")),
		},
		"it should remove zone": {
			zoneID: "1",
			zon:    zon,
			zons:   &zons,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := zone.RepositoryMock{}
			log := logger.LoggerMock{}
			remove := zone.NewRemover(&repo, &log)

			repo.FindFunc = func(id string) *zone.Zone {
				return tt.zon
			}
			log.FatalFunc = func(v ...interface{}) {
			}
			repo.GetZonesFunc = func() *zone.Zones {
				return tt.zons
			}
			repo.SaveFunc = func(z zone.Zones) error {
				return tt.err
			}

			err := remove.Remove(tt.zoneID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
