package execution

import "github.com/bxcodec/faker/v3"

func NewLogsStub() (Logs, error) {
	l := Logs{}
	if err := faker.FakeData(&l); err != nil {
		return Logs{}, err
	}

	return l, nil
}

func NewExecutionStub() Execution {
	programs := ProgramsStub()
	weeklyPrograms := WeeklyProgramsStub()
	exec, _ := New(&programs, &weeklyPrograms, &programs, &programs)
	return *exec
}

func WeeklyProgramsStub() WeeklyPrograms {
	w := WeeklyPrograms{}
	weekly := weeklyStub()
	w.Add(&weekly)
	return w
}

func weeklyStub() Weekly {
	programs := ProgramsStub()
	w := NewWeeklyByDay(&programs, 0)
	return *w
}

func ProgramsStub() Programs {
	p := Programs{}
	program := ProgramStub()
	p.Add(&program)
	return p
}

func ProgramStub() Program {
	p, _ := NewProgram(15, "21:00", []string{"1", "2"})
	return *p
}
