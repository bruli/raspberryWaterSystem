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

func New(seconds Seconds, hour Hour, zones []string) (Program, error) {
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
