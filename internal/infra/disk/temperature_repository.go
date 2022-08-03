package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type temperatureMap = map[float32]programMap

type TemperatureRepository struct {
	path string
}

func (t TemperatureRepository) Save(ctx context.Context, programs []program.Temperature) error {
	temp := make(temperatureMap)
	for _, pr := range programs {
		temp[pr.Temperature()] = buildProgramMap(pr.Programs())
	}
	return writeFile(t.path, temp)
}

func (t TemperatureRepository) FindAll(ctx context.Context) ([]program.Temperature, error) {
	temperature := make(temperatureMap)
	if err := readFile(t.path, &temperature); err != nil {
		return nil, err
	}
	return buildTemperaturePrograms(temperature), nil
}

func buildTemperaturePrograms(temperature temperatureMap) []program.Temperature {
	prgms := make([]program.Temperature, 0, len(temperature))
	for temp, t := range temperature {
		var prg program.Temperature
		prg.Hydrate(temp, buildPrograms(t))
		prgms = append(prgms, prg)
	}
	return prgms
}

func (t TemperatureRepository) FindByTemperatureAndHour(ctx context.Context, temperature float32, hour program.Hour) (program.Temperature, error) {
	temp := make(temperatureMap)
	if err := readFile(t.path, &temp); err != nil {
		return program.Temperature{}, err
	}
	byTemp, ok := temp[temperature]
	if !ok {
		return program.Temperature{}, vo.NotFoundError{}
	}
	byHour, ok := byTemp[hour.String()]
	if !ok {
		return program.Temperature{}, vo.NotFoundError{}
	}
	return buildProgramTemperature(temperature, hour, byHour), nil
}

func buildProgramTemperature(temperature float32, hour program.Hour, prgms []programData) program.Temperature {
	programs := make([]program.Program, 0, len(prgms))
	var temp program.Temperature
	var pg program.Program
	executions := make([]program.Execution, 0, len(prgms))
	for _, pd := range prgms {
		var execution program.Execution
		sec, _ := program.ParseSeconds(pd.Seconds)
		execution.Hydrate(sec, pd.Zones)
		executions = append(executions, execution)
		pg.Hydrate(hour, executions)
		programs = append(programs, pg)
	}
	temp.Hydrate(temperature, programs)
	return temp
}

func NewTemperatureRepository(path string) TemperatureRepository {
	return TemperatureRepository{path: path}
}
