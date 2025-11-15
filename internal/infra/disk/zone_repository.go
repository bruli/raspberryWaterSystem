package disk

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

type (
	zonesMap map[string]zoneData
	zoneData struct {
		Name   string `yaml:"name"`
		Relays []int  `yaml:"relays"`
	}
)

type ZoneRepository struct {
	filePath string
}

func (z ZoneRepository) Update(ctx context.Context, zo *zone.Zone) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		zones := make(zonesMap)
		if err := readYamlFile(z.filePath, &zones); err != nil {
			return err
		}
		zones[zo.Id()] = zoneData{
			Name:   zo.Name(),
			Relays: z.buildRelaysForYaml(zo.Relays()),
		}
		return writeYamlFile(z.filePath, zones)
	}
}

func (z ZoneRepository) FindAll(ctx context.Context) ([]*zone.Zone, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		zones := make(zonesMap)
		if err := readYamlFile(z.filePath, &zones); err != nil {
			return nil, err
		}
		return z.buildZones(zones), nil
	}
}

func (z ZoneRepository) Remove(ctx context.Context, zo *zone.Zone) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		zones := make(zonesMap)
		if err := readYamlFile(z.filePath, &zones); err != nil {
			return err
		}
		_, ok := zones[zo.Id()]
		if !ok {
			return vo.NewNotFoundError(zo.Id())
		}
		delete(zones, zo.Id())
		return writeYamlFile(z.filePath, zones)
	}
}

func (z ZoneRepository) FindByID(ctx context.Context, id string) (*zone.Zone, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		zones := make(zonesMap)
		if err := readYamlFile(z.filePath, &zones); err != nil {
			return nil, err
		}
		zo, ok := zones[id]
		if !ok {
			return nil, vo.NewNotFoundError(id)
		}
		return buildZone(id, zo), nil
	}
}

func buildZone(id string, zo zoneData) *zone.Zone {
	var do zone.Zone
	do.Hydrate(id, zo.Name, buildRelays(zo.Relays))
	return &do
}

func buildRelays(relays []int) []zone.Relay {
	rel := make([]zone.Relay, len(relays))
	for i, n := range relays {
		r, _ := zone.ParseRelay(n)
		rel[i] = r
	}
	return rel
}

func (z ZoneRepository) Save(ctx context.Context, zo *zone.Zone) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		zones := make(zonesMap)
		if err := readYamlFile(z.filePath, &zones); err != nil {
			return err
		}
		zones[zo.Id()] = zoneData{
			Name:   zo.Name(),
			Relays: z.buildRelaysForYaml(zo.Relays()),
		}
		return writeYamlFile(z.filePath, zones)
	}
}

func (z ZoneRepository) buildRelaysForYaml(rel []zone.Relay) []int {
	relays := make([]int, len(rel))
	for i, re := range rel {
		relays[i] = re.Id().Int()
	}
	return relays
}

func (z ZoneRepository) buildZones(data zonesMap) []*zone.Zone {
	zones := make([]*zone.Zone, 0, len(data))
	for i, zo := range data {
		zones = append(zones, buildZone(i, zo))
	}
	return zones
}

func NewZoneRepository(filePath string) ZoneRepository {
	return ZoneRepository{filePath: filePath}
}
