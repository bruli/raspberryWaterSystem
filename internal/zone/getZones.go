package zone

type Getter struct {
	ZoneRepository Repository
}

func NewGetter(zoneRepository Repository) *Getter {
	return &Getter{ZoneRepository: zoneRepository}
}

func (g *Getter) Get() *Zones {
	return g.ZoneRepository.GetZones()
}
