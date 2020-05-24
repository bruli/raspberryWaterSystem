package zone

import error2 "github.com/bruli/raspberryWaterSystem/internal/error"

type Zones []Zone

func (z *Zones) Add(zo Zone) {
	*z = append(*z, zo)
}

func (z *Zones) remove(id string) {
	zs := Zones{}
	for _, j := range *z {
		if id != j.id {
			zs.Add(j)
		}
	}
	*z = zs
}

type Zone struct {
	id     string
	name   string
	relays []string
}

func (z *Zone) Relays() []string {
	return z.relays
}

func (z *Zone) SetRelays(relays []string) {
	z.relays = relays
}

func (z *Zone) Name() string {
	return z.name
}

func (z *Zone) SetName(name string) {
	z.name = name
}

func (z *Zone) Id() string {
	return z.id
}

func (z *Zone) SetId(id string) {
	z.id = id
}

func (z *Zone) update(name string, relays []string) {
	z.name = name
	z.relays = relays
}

func New(id string, name string, relays []string) (*Zone, error) {
	err := error2.Aggregated{}
	if id == "" {
		err.Add("id cannot be empty")
	}

	if name == "" {
		err.Add("name cannot be empty")
	}

	if 0 == len(relays) {
		err.Add("relays cannot be empty")
	}

	if err.WithErrors() {
		return nil, NewCreateError(err.Error())
	}
	return &Zone{id: id, name: name, relays: relays}, nil
}
