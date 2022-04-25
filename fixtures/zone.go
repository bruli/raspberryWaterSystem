package fixtures

import (
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/google/uuid"
)

type ZoneBuilder struct {
	ID, Name *string
	Relays   []string
}

func (b ZoneBuilder) Build() zone.Zone {
	var z zone.Zone
	id := uuid.New().String()
	if b.ID != nil {
		id = *b.ID
	}
	name := "zone name"
	if b.Name != nil {
		name = *b.Name
	}
	relays := []string{
		"1",
		"2",
	}
	if b.Relays != nil {
		relays = b.Relays
	}

	z.Hydrate(id, name, relays)
	return z
}
