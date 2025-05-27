package disk

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type (
	programMap  = map[string][]programData
	programData = struct {
		Seconds int      `yaml:"seconds"`
		Zones   []string `yaml:"zones"`
	}
)

type ProgramRepository struct {
	filePath string
}

func (d ProgramRepository) Save(ctx context.Context, programs []program.Program) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		dailyPrgms := buildProgramMap(programs)
		return writeYamlFile(d.filePath, dailyPrgms)
	}
}

func buildProgramMap(programs []program.Program) programMap {
	dailyPrgms := make(programMap)
	for _, pr := range programs {
		pgrData := make([]programData, len(pr.Executions()))
		for i, p := range pr.Executions() {
			pgrData[i] = programData{
				Seconds: p.Seconds().Int(),
				Zones:   p.Zones(),
			}
		}
		dailyPrgms[pr.Hour().String()] = pgrData
	}
	return dailyPrgms
}

func (d ProgramRepository) FindAll(ctx context.Context) ([]program.Program, error) {
	dailyPgrms := make(programMap)
	if err := readYamlFile(d.filePath, &dailyPgrms); err != nil {
		return nil, err
	}
	return buildDailyPrograms(dailyPgrms), nil
}

func buildDailyPrograms(pr programMap) []program.Program {
	dailies := make([]program.Program, 0, len(pr))
	for hour, pg := range pr {
		dailies = append(dailies, buildProgram(pg, hour))
	}
	return dailies
}

func (d ProgramRepository) FindByHour(ctx context.Context, hour program.Hour) (program.Program, error) {
	dailyPgrms := make(programMap)
	if err := readYamlFile(d.filePath, &dailyPgrms); err != nil {
		return program.Program{}, err
	}
	pgr, ok := dailyPgrms[hour.String()]
	if !ok {
		return program.Program{}, vo.NotFoundError{}
	}
	return buildProgram(pgr, hour.String()), nil
}

func buildProgram(pgr []programData, hour string) program.Program {
	var prog program.Program
	ho, _ := program.ParseHour(hour)
	exec := make([]program.Execution, len(pgr))
	for i, p := range pgr {
		var pExec program.Execution
		sec, _ := program.ParseSeconds(p.Seconds)
		pExec.Hydrate(sec, p.Zones)
		exec[i] = pExec
	}
	prog.Hydrate(ho, exec)
	return prog
}

func NewProgramRepository(filePath string) ProgramRepository {
	return ProgramRepository{filePath: filePath}
}
