package disk

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type weeklyMap = map[string]programMap

type WeeklyRepository struct {
	path string
}

func (w WeeklyRepository) FindByDayAndHour(ctx context.Context, day program.WeekDay, hour program.Hour) (program.Weekly, error) {
	select {
	case <-ctx.Done():
		return program.Weekly{}, ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return program.Weekly{}, err
		}
		byDay, ok := weekly[day.String()]
		if !ok {
			return program.Weekly{}, vo.NotFoundError{}
		}
		byHour, ok := byDay[hour.String()]
		if !ok {
			return program.Weekly{}, vo.NotFoundError{}
		}
		return buildProgramWeekly(day, hour, byHour), nil
	}
}

func buildProgramWeekly(day program.WeekDay, hour program.Hour, prgms []programData) program.Weekly {
	programs := make([]program.Program, 0, len(prgms))
	var weekly program.Weekly
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
	weekly.Hydrate(day, programs)
	return weekly
}

func (w WeeklyRepository) Save(ctx context.Context, programs []program.Weekly) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		weekly := make(weeklyMap)
		for _, pr := range programs {
			weekly[pr.WeekDay().String()] = buildProgramMap(pr.Programs())
		}
		return writeYamlFile(w.path, weekly)
	}
}

func (w WeeklyRepository) FindAll(ctx context.Context) ([]program.Weekly, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return nil, err
		}
		return buildWeeklyPrograms(weekly), nil
	}
}

func buildWeeklyPrograms(weekly weeklyMap) []program.Weekly {
	prgms := make([]program.Weekly, 0, len(weekly))
	for dayStr, w := range weekly {
		day, _ := program.ParseWeekDay(dayStr)
		var prg program.Weekly
		prg.Hydrate(day, buildPrograms(w))
		prgms = append(prgms, prg)
	}
	return prgms
}

func buildPrograms(w programMap) []program.Program {
	pgs := make([]program.Program, 0, len(w))
	for hour, pg := range w {
		pgs = append(pgs, buildProgram(pg, hour))
	}
	return pgs
}

func NewWeeklyRepository(path string) WeeklyRepository {
	return WeeklyRepository{path: path}
}
