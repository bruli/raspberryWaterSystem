package disk

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"gopkg.in/yaml.v3"
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

func (z ZoneRepository) Update(ctx context.Context, zo zone.Zone) error {
	// TODO implement me
	panic("implement me")
}

func NewZoneRepository(filePath string) ZoneRepository {
	return ZoneRepository{filePath: filePath}
}

func (z ZoneRepository) FindByID(_ context.Context, id string) (zone.Zone, error) {
	zones := make(zonesMap)
	if err := readFile(z.filePath, &zones); err != nil {
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

func readFile(path string, data interface{}) error {
	if err := checkFile(path); err != nil {
		return fmt.Errorf("failed to check %s file", path)
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if err = yaml.Unmarshal(fileData, data); err != nil {
		return err
	}
	return nil
}

func checkFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = ioutil.WriteFile(path, nil, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

func (z ZoneRepository) Save(_ context.Context, zo zone.Zone) error {
	zones := make(zonesMap)
	if err := readFile(z.filePath, &zones); err != nil {
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
	return writeFile(z.filePath, zones)
}

func writeFile(path string, data interface{}) error {
	dataFile, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return ioutil.WriteFile(path, dataFile, 0o755)
}
