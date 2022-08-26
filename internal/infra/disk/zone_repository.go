package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
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

func (z ZoneRepository) Remove(_ context.Context, zo zone.Zone) error {
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

func NewZoneRepository(filePath string) ZoneRepository {
	return ZoneRepository{filePath: filePath}
}

func (z ZoneRepository) FindByID(_ context.Context, id string) (zone.Zone, error) {
	zones := make(zonesMap)
	if err := readYamlFile(z.filePath, &zones); err != nil {
		return zone.Zone{}, err
	}
	zo, ok := zones[id]
	if !ok {
		return zone.Zone{}, vo.NewNotFoundError(id)
	}
	return buildZone(id, zo), nil
}

func buildZone(id string, zo zoneData) zone.Zone {
	var do zone.Zone
	do.Hydrate(id, zo.Name, buildRelays(zo.Relays))
	return do
}

func buildRelays(relays []int) []zone.Relay {
	rel := make([]zone.Relay, len(relays))
	for i, n := range relays {
		r, _ := zone.ParseRelay(n)
		rel[i] = r
	}
	return rel
}

func (z ZoneRepository) Save(_ context.Context, zo zone.Zone) error {
	zones := make(zonesMap)
	if err := readYamlFile(z.filePath, &zones); err != nil {
		return err
	}
	relays := make([]int, len(zo.Relays()))
	for i, re := range zo.Relays() {
		relays[i] = re.Id().Int()
	}
	zones[zo.Id()] = zoneData{
		Name:   zo.Name(),
		Relays: relays,
	}
	return writeYamlFile(z.filePath, zones)
}
