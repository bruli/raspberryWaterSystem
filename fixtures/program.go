package fixtures

import "github.com/bruli/raspberryWaterSystem/internal/domain/program"

type ProgramBuilder struct {
	Seconds *program.Seconds
	Hour    *program.Hour
	Zones   []string
}

func (b ProgramBuilder) Build() program.Program {
	var pr program.Program
	seconds, _ := program.ParseSeconds(20)
	if b.Seconds != nil {
		seconds = *b.Seconds
	}
	hour, _ := program.ParseHour("15:10")
	if b.Hour != nil {
		hour = *b.Hour
	}
	pr.Hydrate(seconds, hour, b.Zones)
	return pr
}

type DailyProgramBuilder struct {
	Seconds *program.Seconds
	Hour    *program.Hour
	Zones   []string
}

func (b DailyProgramBuilder) Build() program.Daily {
	var pr program.Daily
	seconds, _ := program.ParseSeconds(20)
	if b.Seconds != nil {
		seconds = *b.Seconds
	}
	hour, _ := program.ParseHour("15:10")
	if b.Hour != nil {
		hour = *b.Hour
	}
	pr.Hydrate(seconds, hour, b.Zones)
	return pr
}

type OddProgramBuilder struct {
	Seconds *program.Seconds
	Hour    *program.Hour
	Zones   []string
}

func (b OddProgramBuilder) Build() program.Odd {
	var pr program.Odd
	seconds, _ := program.ParseSeconds(20)
	if b.Seconds != nil {
		seconds = *b.Seconds
	}
	hour, _ := program.ParseHour("15:10")
	if b.Hour != nil {
		hour = *b.Hour
	}
	pr.Hydrate(seconds, hour, b.Zones)
	return pr
}

type EvenProgramBuilder struct {
	Seconds *program.Seconds
	Hour    *program.Hour
	Zones   []string
}

func (b EvenProgramBuilder) Build() program.Even {
	var pr program.Even
	seconds, _ := program.ParseSeconds(20)
	if b.Seconds != nil {
		seconds = *b.Seconds
	}
	hour, _ := program.ParseHour("15:10")
	if b.Hour != nil {
		hour = *b.Hour
	}
	pr.Hydrate(seconds, hour, b.Zones)
	return pr
}
