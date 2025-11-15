package fixtures

import (
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/google/uuid"
)

type ZoneBuilder struct {
	ID, Name *string
	Relays   []zone.Relay
}

func (b ZoneBuilder) Build() *zone.Zone {
	var z zone.Zone
	id := uuid.NewString()
	if b.ID != nil {
		id = *b.ID
	}
	name := "zone name"
	if b.Name != nil {
		name = *b.Name
	}
	relays := []zone.Relay{
		RelayBuilder(1),
		RelayBuilder(2),
		RelayBuilder(3),
	}
	if b.Relays != nil {
		relays = b.Relays
	}

	z.Hydrate(id, name, relays)
	return &z
}

func RelayBuilder(i int) zone.Relay {
	rel, _ := zone.ParseRelay(i)
	return rel
}
