package zone_test

import (
	zone2 "github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	tests := map[string]struct {
		id, name string
		relays   []string
		err      error
	}{
		"it should return errors with invalid data": {
			id:     "",
			name:   "",
			relays: []string{},
			err:    zone2.NewCreateError("id cannot be empty,name cannot be empty,relays cannot be empty"),
		},
		"it should return zone": {
			id:     "aaa",
			name:   "esto",
			relays: []string{"1"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			z, err := zone2.New(tt.id, tt.name, tt.relays)
			assert.Equal(t, tt.err, err)
			if z != nil {
				assert.Equal(t, tt.id, z.Id())
				assert.Equal(t, tt.name, z.Name())
				assert.Equal(t, tt.relays, z.Relays())
			}
		})
	}
}
