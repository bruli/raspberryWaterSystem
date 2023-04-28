package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type temperatureMap = map[float32]programMap

type TemperatureProgramRepository struct {
	path string
}

func (t TemperatureProgramRepository) Save(ctx context.Context, programs []program.Temperature) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		temp := make(temperatureMap)
		for _, pr := range programs {
			temp[pr.Temperature()] = buildProgramMap(pr.Programs())
		}
		return writeYamlFile(t.path, temp)
	}
}

func (t TemperatureProgramRepository) FindAll(ctx context.Context) ([]program.Temperature, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		temperature := make(temperatureMap)
		if err := readYamlFile(t.path, &temperature); err != nil {
			return nil, err
		}
		return buildTemperaturePrograms(temperature), nil
	}
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

func (t TemperatureProgramRepository) FindByTemperatureAndHour(ctx context.Context, temperature float32, hour program.Hour) (program.Temperature, error) {
	select {
	case <-ctx.Done():
		return program.Temperature{}, ctx.Err()
	default:
		temp := make(temperatureMap)
		if err := readYamlFile(t.path, &temp); err != nil {
			return program.Temperature{}, err
		}
		programs := make(programMap, 0)
		for tempKey, progs := range temp {
			if temperature >= tempKey {
				for hourKey, prgms := range progs {
					programs[hourKey] = prgms
				}
			}
		}
		byHour, ok := programs[hour.String()]
		if !ok {
			return program.Temperature{}, vo.NotFoundError{}
		}
		return buildProgramTemperature(temperature, hour, byHour), nil
	}
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

func NewTemperatureProgramRepository(path string) TemperatureProgramRepository {
	return TemperatureProgramRepository{path: path}
}
