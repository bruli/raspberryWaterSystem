package execution

type TemperaturePrograms []*TemperatureProgram

func (t *TemperaturePrograms) add(p *TemperatureProgram) {
	*t = append(*t, p)
}

func (t *TemperaturePrograms) GetPrograms(temp float32) *Programs {
	pgrms := Programs{}
	for _, pgr := range *t {
		if temp >= pgr.Temperature {
			pgrms.Add(&pgr.Program)
		}
	}
	return &pgrms
}

func (t *TemperaturePrograms) Add(program *TemperatureProgram) {
	*t = append(*t, program)
}
