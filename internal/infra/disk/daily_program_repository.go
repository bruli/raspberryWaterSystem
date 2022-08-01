package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type (
	programMap  = map[string]programData
	programData = struct {
		Seconds int      `yaml:"seconds"`
		Zones   []string `yaml:"zones"`
	}
)

type DailyProgramRepository struct {
	filePath string
}

func (d DailyProgramRepository) Save(ctx context.Context, programs []program.Daily) error {
	dailyPrgms := make(programMap)
	for _, pr := range programs {
		dailyPrgms[pr.Hour().String()] = programData{
			Seconds: pr.Seconds().Int(),
			Zones:   pr.Zones(),
		}
	}
	return writeFile(d.filePath, dailyPrgms)
}

func (d DailyProgramRepository) FindAll(ctx context.Context) ([]program.Daily, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return nil, err
	}
	return buildDailyPrograms(dailyPgrms), nil
}

func buildDailyPrograms(pr programMap) []program.Daily {
	dailies := make([]program.Daily, 0, len(pr))
	for hour, pg := range pr {
		dailies = append(dailies, buildDaily(pg, hour))
	}
	return dailies
}

func (d DailyProgramRepository) FindByHour(ctx context.Context, hour program.Hour) (program.Daily, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return program.Daily{}, err
	}
	pgr, ok := dailyPgrms[hour.String()]
	if !ok {
		return program.Daily{}, vo.NotFoundError{}
	}
	return buildDaily(pgr, hour.String()), nil
}

func buildDaily(pgr programData, hour string) program.Daily {
	var prog program.Daily
	sec, _ := program.ParseSeconds(pgr.Seconds)
	ho, _ := program.ParseHour(hour)
	prog.Hydrate(sec, ho, pgr.Zones)
	return prog
}

func NewDailyProgramRepository(filePath string) DailyProgramRepository {
	return DailyProgramRepository{filePath: filePath}
}
