package disk

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type (
	programMap = map[string][]executions
	executions = struct {
		Seconds int      `yaml:"seconds"`
		Zones   []string `yaml:"zones"`
	}
)

type ProgramRepository struct {
	filePath string
}

func (d ProgramRepository) Save(ctx context.Context, prg *program.Program) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		savedData := make(programMap)
		if err := readYamlFile(d.filePath, &savedData); err != nil {
			return err
		}
		execs := make([]executions, len(prg.Executions()))
		for i, p := range prg.Executions() {
			execs[i] = buildExecution(&p)
		}
		savedData[prg.Hour().String()] = execs
		return writeYamlFile(d.filePath, savedData)
	}
}

func buildExecution(exec *program.Execution) executions {
	return executions{
		Seconds: exec.Seconds().Int(),
		Zones:   exec.Zones(),
	}
}

func buildProgramMap(programs []program.Program) programMap {
	dailyPrgms := make(programMap)
	for _, pr := range programs {
		pgrData := make([]executions, len(pr.Executions()))
		for i, p := range pr.Executions() {
			pgrData[i] = executions{
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
		p := buildProgram(pg, hour)
		dailies = append(dailies, *p)
	}
	return dailies
}

func (d ProgramRepository) FindByHour(ctx context.Context, hour *program.Hour) (*program.Program, error) {
	dailyPgrms := make(programMap)
	if err := readYamlFile(d.filePath, &dailyPgrms); err != nil {
		return nil, err
	}
	pgr, ok := dailyPgrms[hour.String()]
	if !ok {
		return nil, vo.NewNotFoundError(hour.String())
	}
	return buildProgram(pgr, hour.String()), nil
}

func buildProgram(pgr []executions, hour string) *program.Program {
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
	return &prog
}

func NewProgramRepository(filePath string) ProgramRepository {
	return ProgramRepository{filePath: filePath}
}
