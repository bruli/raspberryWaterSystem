package program

import (
	"errors"
)

var (
	ErrEmptyPrograms          = errors.New("empty programs")
	ErrEmptyExecutionZones    = errors.New("empty execution zones")
	ErrEmptyProgramExecutions = errors.New("empty program executions")
)

type Program struct {
	hour       Hour
	executions []Execution
}

func (p Program) Hour() Hour {
	return p.hour
}

func (p Program) Executions() []Execution {
	return p.executions
}

func (p Program) validate() error {
	if len(p.executions) == 0 {
		return ErrEmptyProgramExecutions
	}
	return nil
}

func (p *Program) Hydrate(hour Hour, executions []Execution) {
	p.hour = hour
	p.executions = executions
}

func New(hour Hour, executions []Execution) (Program, error) {
	pr := Program{hour: hour, executions: executions}
	if err := pr.validate(); err != nil {
		return Program{}, err
	}
	return pr, nil
}

type Weekly struct {
	weekDay  WeekDay
	programs []Program
}

func (w Weekly) WeekDay() WeekDay {
	return w.weekDay
}

func (w Weekly) Programs() []Program {
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
