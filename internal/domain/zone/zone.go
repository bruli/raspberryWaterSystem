package zone

import "errors"

var (
	ErrInvalidZoneID     = errors.New("zone id can not be empty")
	ErrInvalidZoneName   = errors.New("zone name can not be empty")
	ErrInvalidZoneRelays = errors.New("zone relays can not be empty")
)

type Zone struct {
	id, name string
	relays   []string
}

func (z Zone) Id() string {
	return z.id
}

func (z Zone) Name() string {
	return z.name
}

func (z Zone) Relays() []string {
	return z.relays
}

func New(id, name string, relays []string) (Zone, error) {
	z := Zone{id: id, name: name, relays: relays}
	if err := z.validate(); err != nil {
		return Zone{}, err
	}
	return z, nil
}

func (z *Zone) Hydrate(id, name string, relays []string) {
	z.id = id
	z.name = name
	z.relays = relays
}

func (z Zone) validate() error {
	if len(z.id) == 0 {
		return ErrInvalidZoneID
	}
	if len(z.name) == 0 {
		return ErrInvalidZoneName
	}
	if len(z.relays) == 0 {
		return ErrInvalidZoneRelays
	}
	return nil
}

func (z *Zone) Update(name string, relays []string) {
	z.name = name
	z.relays = relays
}