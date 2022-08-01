package program

type Program struct {
	seconds   Seconds
	execution Execution
}

func (p Program) Seconds() Seconds {
	return p.seconds
}

func (p Program) Execution() Execution {
	return p.execution
}

func new(seconds Seconds, hour Hour, zones []string) (Program, error) {
	exec, err := NewExecution(hour, zones)
	if err != nil {
		return Program{}, err
	}
	if _, err = ParseSeconds(seconds.Int()); err != nil {
		return Program{}, err
	}
	return Program{
		seconds:   seconds,
		execution: exec,
	}, nil
}

type DailyProgram struct {
	Program
}

func NewDaily(seconds Seconds, hour Hour, zones []string) (DailyProgram, error) {
	pr, err := new(seconds, hour, zones)
	if err != nil {
		return DailyProgram{}, err
	}
	return DailyProgram{Program: pr}, nil
}

type OddProgram struct {
	Program
}

func NewOdd(seconds Seconds, hour Hour, zones []string) (OddProgram, error) {
	pr, err := new(seconds, hour, zones)
	if err != nil {
		return OddProgram{}, err
	}
	return OddProgram{Program: pr}, nil
}

type EvenProgram struct {
	Program
}

func NewEven(seconds Seconds, hour Hour, zones []string) (EvenProgram, error) {
	pr, err := new(seconds, hour, zones)
	if err != nil {
		return EvenProgram{}, err
	}
	return EvenProgram{Program: pr}, nil
}
