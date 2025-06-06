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

func (w WeeklyRepository) FindByDay(ctx context.Context, day *program.WeekDay) (*program.Weekly, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return nil, err
		}
		byDay, ok := weekly[day.String()]
		if !ok {
			return nil, vo.NotFoundError{}
		}
		return buildProgramWeeklyByDay(day, byDay), nil
	}
}

func buildProgramWeeklyByDay(day *program.WeekDay, prg programMap) *program.Weekly {
	var weekly program.Weekly

	weekly.Hydrate(*day, buildPrograms(prg))
	return &weekly
}

func (w WeeklyRepository) Remove(ctx context.Context, day *program.WeekDay) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return err
		}
		delete(weekly, day.String())
		if err := writeYamlFile(w.path, &weekly); err != nil {
			return err
		}
		return nil
	}
}

func (w WeeklyRepository) FindByDayAndHour(ctx context.Context, day *program.WeekDay, hour *program.Hour) (*program.Weekly, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return nil, err
		}
		byDay, ok := weekly[day.String()]
		if !ok {
			return nil, vo.NewNotFoundError(day.String())
		}
		byHour, ok := byDay[hour.String()]
		if !ok {
			return nil, vo.NewNotFoundError(hour.String())
		}
		return buildProgramWeeklyByHour(day, hour, byHour), nil
	}
}

func buildProgramWeeklyByHour(day *program.WeekDay, hour *program.Hour, prgms []executions) *program.Weekly {
	programs := make([]program.Program, 0, len(prgms))
	var weekly program.Weekly
	var pg program.Program
	exec := make([]program.Execution, 0, len(prgms))
	for _, pd := range prgms {
		var execution program.Execution
		sec, _ := program.ParseSeconds(pd.Seconds)
		execution.Hydrate(sec, pd.Zones)
		exec = append(exec, execution)
		pg.Hydrate(*hour, exec)
		programs = append(programs, pg)
	}
	weekly.Hydrate(*day, programs)
	return &weekly
}

func (w WeeklyRepository) Save(ctx context.Context, program *program.Weekly) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		weekly := make(weeklyMap)
		if err := readYamlFile(w.path, &weekly); err != nil {
			return err
		}
		weekly[program.WeekDay().String()] = buildProgramMap(program.Programs())
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

func NewWeeklyRepository(path string) WeeklyRepository {
	return WeeklyRepository{path: path}
}
