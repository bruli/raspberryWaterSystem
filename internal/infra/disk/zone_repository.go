package disk

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/davecgh/go-spew/spew"
)

type zoneData struct {
	ID     string   `yaml:"id"`
	Name   string   `yaml:"name"`
	Relays []string `yaml:"relays"`
}

type ZoneRepository struct {
	filePath string
}

func NewZoneRepository(filePath string) ZoneRepository {
	return ZoneRepository{filePath: filePath}
}

func (z ZoneRepository) FindByID(ctx context.Context, id string) (zone.Zone, error) {
	zones, err := readFile(z.filePath)
	if err != nil {
		return zone.Zone{}, err
	}
	zo, ok := zones[id]
	if !ok {
		return zone.Zone{}, vo.NewNotFoundError(id)
	}
	var zon zone.Zone
	zon.Hydrate(zo.ID, zo.Name, zo.Relays)
	return zon, nil
}

func readFile(path string) (map[string]zoneData, error) {
	if err := checkFile(path); err != nil {
		return nil, err
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	var zonesData []zoneData
	if err := yaml.Unmarshal(fileData, &zonesData); err != nil {
		return nil, err
	}
	zonesMap := make(map[string]zoneData, len(zonesData))
	for _, zo := range zonesData {
		zonesMap[zo.ID] = zo
	}
	return zonesMap, nil
}

func checkFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := ioutil.WriteFile(path, nil, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

func (z ZoneRepository) Save(_ context.Context, zo zone.Zone) error {
	zones, err := readFile(z.filePath)
	if err != nil {
		return err
	}
	_, ok := zones[zo.Id()]
	if ok {
		return DuplicatedZoneError{id: zo.Id()}
	}
	zones[zo.Id()] = buildZoneData(zo)
	return writeFile(z.filePath, zones)
}

func writeFile(path string, zones map[string]zoneData) error {
	var zonesData []zoneData
	for _, zo := range zones {
		zonesData = append(zonesData, zo)
	}
	dataFile, err := yaml.Marshal(zonesData)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return ioutil.WriteFile(path, dataFile, 0755)
}

func buildZoneData(zo zone.Zone) zoneData {
	return zoneData{
		ID:     zo.Id(),
		Name:   zo.Name(),
		Relays: zo.Relays(),
	}
}

func (z ZoneRepository) Update(_ context.Context, zo zone.Zone) error {
	zones, err := readFile(z.filePath)
	if err != nil {
		return err
	}
	_, ok := zones[zo.Id()]
	if !ok {
		return vo.NewNotFoundError(zo.Id())
	}
	zones[zo.Id()] = buildZoneData(zo)
	return writeFile(z.filePath, zones)
}

type DuplicatedZoneError struct {
	id string
}

func (d DuplicatedZoneError) Error() string {
	return spew.Sprintf("duplicated zone: %s", d.id)
}
