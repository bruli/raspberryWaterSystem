package zone

import (
	"github.com/bxcodec/faker/v3"
)

type ZoneStub struct {
	ID     string
	Name   string
	Relays []string
}

func NewZoneStub() (Zone, error) {
	z := ZoneStub{}
	_ = faker.SetRandomMapAndSliceSize(3)
	err := faker.FakeData(&z)
	if err != nil {
		return Zone{}, err
	}
	relays := z.Relays
	if 0 == len(relays) {
		relays = []string{"1"}
	}
	zo, err := New(z.ID, z.Name, relays)
	if err != nil {
		return Zone{}, err
	}
	return *zo, nil
}

func NewZonesStub() (Zones, error) {
	z := Zones{}
	zo, err := NewZoneStub()
	if err != nil {
		return z, err
	}
	z.Add(zo)

	return z, nil
}
