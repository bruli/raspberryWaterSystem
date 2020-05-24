package execution

import "time"

type Program struct {
	Seconds    time.Duration
	Executions *Data
}

func NewProgram(seconds uint8, hour string, zones []string) (*Program, error) {
	exT, _ := time.Parse("15:04", hour)
	data, err := NewData(exT, zones)
	if err != nil {
		return nil, err
	}
	return &Program{Seconds: time.Duration(seconds) * time.Second, Executions: data}, nil
}

func (p *Program) getHour() string {
	return p.Executions.Hour.Format("15:04")
}

func (p *Program) hasSameZones(prgm *Program) bool {
	for i, zon := range p.Executions.Zones {
		if zon != prgm.Executions.Zones[i] {
			return false
		}
	}

	return true
}

func (c *Program) getBestProgram(prgm *Program) *Program {
	if c.hasSameZones(prgm) {
		if prgm.Seconds > c.Seconds {
			return prgm
		}
	}

	return c
}
