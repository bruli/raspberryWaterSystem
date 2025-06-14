package disk

import (
	"context"
	"strconv"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type temperatureMap = map[float32]programMap

type TemperatureProgramRepository struct {
	path string
}

func (t TemperatureProgramRepository) FindByTemperature(ctx context.Context, temperature float32) (*program.Temperature, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		temp := make(temperatureMap)
		if err := readYamlFile(t.path, &temp); err != nil {
			return nil, err
		}
		byTemp, ok := temp[temperature]
		if !ok {
			return nil, vo.NewNotFoundError(strconv.FormatFloat(float64(temperature), 'f', -1, 32))
		}
		return buildTemperatureProgram(temperature, byTemp), nil
	}
}

func buildTemperatureProgram(temperature float32, prgms programMap) *program.Temperature {
	var tempPrgm program.Temperature
	tempPrgm.Hydrate(temperature, buildPrograms(prgms))
	return &tempPrgm
}

func (t TemperatureProgramRepository) Remove(ctx context.Context, temperature float32) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		temp := make(temperatureMap)
		if err := readYamlFile(t.path, &temp); err != nil {
			return err
		}
		delete(temp, temperature)
		if err := writeYamlFile(t.path, &temp); err != nil {
			return err
		}
		return nil
	}
}

func (t TemperatureProgramRepository) Save(ctx context.Context, programs *program.Temperature) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		temp := make(temperatureMap)
		if err := readYamlFile(t.path, &temp); err != nil {
			return err
		}
		temp[programs.Temperature()] = buildProgramMap(programs.Programs())
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

func buildProgramTemperature(temperature float32, hour program.Hour, prgms []executions) program.Temperature {
	programs := make([]program.Program, len(prgms))
	var temp program.Temperature
	var pg program.Program
	for i, pd := range prgms {
		pg.Hydrate(hour, buildExecutions(pd))
		programs[i] = pg
	}
	temp.Hydrate(temperature, programs)
	return temp
}

func buildExecutions(pd executions) []program.Execution {
	executions := make([]program.Execution, 1)
	var execution program.Execution
	sec, _ := program.ParseSeconds(pd.Seconds)
	execution.Hydrate(sec, pd.Zones)
	executions[0] = execution
	return executions
}

func NewTemperatureProgramRepository(path string) TemperatureProgramRepository {
	return TemperatureProgramRepository{path: path}
}
