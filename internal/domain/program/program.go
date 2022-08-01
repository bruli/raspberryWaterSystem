package program

import (
	"errors"
	"time"
)

var (
	ErrEmptyWeeklyPrograms = errors.New("empty weekly programs")
	ErrEmptyExecutionZones = errors.New("empty execution zones")
)

type Program struct {
	seconds Seconds
	hour    Hour
	zones   []string
}

func (p Program) Hour() Hour {
	return p.hour
}

func (p Program) Zones() []string {
	return p.zones
}

func (p Program) Seconds() Seconds {
	return p.seconds
}

func program(seconds Seconds, hour Hour, zones []string) (Program, error) {
	pr := Program{
		seconds: seconds,
		hour:    hour,
		zones:   zones,
	}
	if err := pr.validate(); err != nil {
		return Program{}, err
	}
	return pr, nil
}

func (p *Program) Hydrate(seconds Seconds, hour Hour, zones []string) {
	p.seconds = seconds
	p.hour = hour
	p.zones = zones
}

func (p Program) validate() error {
	if len(p.zones) == 0 {
		return ErrEmptyExecutionZones
	}
	if _, err := ParseSeconds(p.seconds.Int()); err != nil {
		return err
	}
	return nil
}

type Daily struct {
	Program
}

func NewDaily(seconds Seconds, hour Hour, zones []string) (Daily, error) {
	pr, err := program(seconds, hour, zones)
	if err != nil {
		return Daily{}, err
	}
	return Daily{Program: pr}, nil
}

func (d *Daily) Hydrate(seconds Seconds, hour Hour, zones []string) {
	d.hour = hour
	d.seconds = seconds
	d.zones = zones
}

type Odd struct {
	Program
}

func NewOdd(seconds Seconds, hour Hour, zones []string) (Odd, error) {
	pr, err := program(seconds, hour, zones)
	if err != nil {
		return Odd{}, err
	}
	return Odd{Program: pr}, nil
}

type Even struct {
	Program
}

func NewEven(seconds Seconds, hour Hour, zones []string) (Even, error) {
	pr, err := program(seconds, hour, zones)
	if err != nil {
		return Even{}, err
	}
	return Even{Program: pr}, nil
}

type Weekly struct {
	weekDay  time.Weekday
	programs []Program
}

func (w Weekly) WeekDay() time.Weekday {
	return w.weekDay
}

func (w Weekly) Programs() []Program {
	return w.programs
}

func NewWeekly(weekDay time.Weekday, programs []Program) (Weekly, error) {
	if len(programs) == 0 {
		return Weekly{}, ErrEmptyWeeklyPrograms
	}
	return Weekly{weekDay: weekDay, programs: programs}, nil
}
