package program

import (
	"errors"
)

var (
	ErrEmptyPrograms       = errors.New("empty programs")
	ErrEmptyExecutionZones = errors.New("empty execution zones")
	ErrInvalidTemperature  = errors.New("invalid temperature")
)

type Program struct {
	hour       Hour
	executions []Execution
}

func (p *Program) Hour() Hour {
	return p.hour
}

func (p *Program) Executions() []Execution {
	return p.executions
}

func (p *Program) validate() error {
	if len(p.executions) == 0 {
		return ErrEmptyPrograms
	}
	return nil
}

func (p *Program) Hydrate(hour Hour, executions []Execution) {
	p.hour = hour
	p.executions = executions
}

func New(hour Hour, executions []Execution) (*Program, error) {
	pr := Program{hour: hour, executions: executions}
	if err := pr.validate(); err != nil {
		return nil, err
	}
	return &pr, nil
}

type Weekly struct {
	weekDay  WeekDay
	programs []Program
}

func (w *Weekly) WeekDay() WeekDay {
	return w.weekDay
}

func (w *Weekly) Programs() []Program {
	return w.programs
}

func NewWeekly(weekDay WeekDay, programs []Program) (Weekly, error) {
	if len(programs) == 0 {
		return Weekly{}, ErrEmptyPrograms
	}
	return Weekly{weekDay: weekDay, programs: programs}, nil
}

func (w *Weekly) Hydrate(weekDay WeekDay, programs []Program) {
	w.weekDay = weekDay
	w.programs = programs
}

type Temperature struct {
	temperature float32
	programs    []Program
}

func (t *Temperature) Temperature() float32 {
	return t.temperature
}

func (t *Temperature) Programs() []Program {
	return t.programs
}

func (t *Temperature) validate() error {
	if len(t.programs) == 0 {
		return ErrEmptyPrograms
	}
	if 0 > t.temperature || 50 < t.temperature {
		return ErrInvalidTemperature
	}
	return nil
}

func NewTemperature(temperature float32, programs []Program) (*Temperature, error) {
	temp := Temperature{temperature: temperature, programs: programs}
	if err := temp.validate(); err != nil {
		return nil, err
	}
	return &temp, nil
}

func (t *Temperature) Hydrate(temperature float32, programs []Program) {
	t.temperature = temperature
	t.programs = programs
}
