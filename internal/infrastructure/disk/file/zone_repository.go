package file

import (
	"fmt"
	zone2 "github.com/bruli/raspberryWaterSystem/internal/zone"
	"gopkg.in/yaml.v2"
)

type zones []*zone

func (z *zones) add(zo *zone) {
	*z = append(*z, zo)
}

type zone struct {
	ID     string   `yaml:"id"`
	Name   string   `yaml:"name"`
	Relays []string `yaml:"relays"`
}

func newZone(ID string, name string, relays []string) *zone {
	return &zone{ID: ID, Name: name, Relays: relays}
}

type ZoneRepository struct {
	repository *repository
}

func NewZoneRepository(file string) *ZoneRepository {
	return &ZoneRepository{repository: newRepository(file)}
}

func (r *ZoneRepository) GetZones() *zone2.Zones {
	data, err := r.repository.reader.read()
	if err != nil {
		return &zone2.Zones{}
	}
	zones := zones{}
	if err := yaml.Unmarshal(data, &zones); err != nil {
		return &zone2.Zones{}
	}

	return r.buildZones(&zones)
}

func (r *ZoneRepository) Save(z zone2.Zones) error {
	data, err := yaml.Marshal(r.buildYamlZones(&z))
	if err != nil {
		return fmt.Errorf("failed to marshal yaml zones: %w", err)
	}
	return r.repository.writer.write(data)
}

func (r *ZoneRepository) Find(id string) *zone2.Zone {
	for _, z := range *r.GetZones() {
		if id == z.Id() {
			return &z
		}
	}

	return nil
}

func (r *ZoneRepository) buildZones(z *zones) *zone2.Zones {
	zo := zone2.Zones{}
	for _, j := range *z {
		zo.Add(*r.buildZone(j))
	}

	return &zo
}

func (r *ZoneRepository) buildZone(j *zone) *zone2.Zone {
	z, err := zone2.New(j.ID, j.Name, j.Relays)
	if err != nil {
		return nil
	}

	return z
}

func (r *ZoneRepository) buildYamlZones(z *zone2.Zones) *zones {
	zon := zones{}
	for _, j := range *z {
		yamlZone := r.buildYamlZone(&j)
		zon.add(yamlZone)
	}
	return &zon
}

func (r *ZoneRepository) buildYamlZone(j *zone2.Zone) *zone {
	return newZone(j.Id(), j.Name(), j.Relays())
}
