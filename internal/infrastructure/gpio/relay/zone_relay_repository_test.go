package relay_test

import (
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/gpio/relay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZoneRelayRepository_Get(t *testing.T) {
	tests := map[string]struct {
		relays []string
	}{
		"it should return relays": {relays: []string{"1", "2", "3", "4", "5", "6"}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			zoneRepo := relay.NewZoneRelayRepository()
			rlys := zoneRepo.Get()
			assert.Equal(t, len(tt.relays), len(rlys))
			for _, r := range tt.relays {
				assert.True(t, func(rel string) bool {
					for _, re := range rlys {
						if re == rel {
							return true
						}
					}
					return false
				}(r))
			}
		})
	}
}
