package program

import (
	"errors"
)

var ErrEmptyExecutionZones = errors.New("empty execution zones")

type Execution struct {
	hour  Hour
	zones []string
}

func (e Execution) Hour() Hour {
	return e.hour
}

func (e Execution) Zones() []string {
	return e.zones
}

func (e Execution) validate() error {
	if len(e.zones) == 0 {
		return ErrEmptyExecutionZones
	}
	return nil
}

func NewExecution(hour Hour, zones []string) (Execution, error) {
	ex := Execution{hour: hour, zones: zones}
	if err := ex.validate(); err != nil {
		return Execution{}, err
	}
	return ex, nil
}
