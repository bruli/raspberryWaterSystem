package zone_test

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	z, _ := zone.New("cc", "aa", []string{"1"})
	tests := map[string]struct {
		id, name                string
		relays, currentRelays   []string
		relaysErr, repoErr, err error
		currentZones            *zone.Zones
		currentZone             *zone.Zone
	}{
		"it should return error with invalid relays": {
			id:            "aaaa",
			name:          "bbbb",
			relays:        []string{"1"},
			currentRelays: []string{"2"},
			err:           zone.NewInvalidRelay("1"),
		},
		"it should return error with invalid zone data": {
			id:            "",
			name:          "bbbb",
			relays:        []string{"1"},
			currentRelays: []string{"1"},
			err:           zone.NewCreateError("id cannot be empty"),
		},
		"it should return error when repository returns error": {
			id:            "aa",
			name:          "bbbb",
			relays:        []string{"1"},
			currentRelays: []string{"1"},
			currentZones:  &zone.Zones{},
			repoErr:       errors.New("error"),
			err:           errors.New("failed saving zones: error"),
		},
		"it should return save new zone": {
			id:            "cc",
			name:          "bbbb",
			relays:        []string{"1"},
			currentRelays: []string{"1"},
			currentZones:  &zone.Zones{},
		},
		"it should return update zone": {
			id:            z.Id(),
			name:          "bbbb",
			relays:        z.Relays(),
			currentRelays: []string{"1"},
			currentZones:  &zone.Zones{},
			currentZone:   z,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := zone.RepositoryMock{}
			relayRepo := zone.RelayRepositoryMock{}
			logger := logger.LoggerMock{}
			create := zone.NewCreator(&repo, &relayRepo, &logger)

			relayRepo.GetFunc = func() []string {
				return tt.currentRelays
			}
			logger.FatalFunc = func(v ...interface{}) {
			}
			repo.GetZonesFunc = func() *zone.Zones {
				return tt.currentZones
			}
			repo.FindFunc = func(id string) *zone.Zone {
				return tt.currentZone
			}
			repo.SaveFunc = func(z zone.Zones) error {
				return tt.repoErr
			}
			err := create.Create(tt.id, tt.name, tt.relays)

			if err != nil {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
