package execution

type TemperatureProgram struct {
	Program
	Temperature float32
}

func NewTemperatureProgram(temperature float32,
	seconds uint8,
	hour string,
	zones []string,
) (*TemperatureProgram, error) {
	p, err := NewProgram(seconds, hour, zones)
	if err != nil {
		return nil, err
	}
	return &TemperatureProgram{Temperature: temperature, Program: *p}, nil
}
