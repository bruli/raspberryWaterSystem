package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type EvenProgramRepository struct {
	filePath string
}

func (d EvenProgramRepository) Save(ctx context.Context, programs []program.Even) error {
	dailyPrgms := make(programMap)
	for _, pr := range programs {
		dailyPrgms[pr.Hour().String()] = programData{
			Seconds: pr.Seconds().Int(),
			Zones:   pr.Zones(),
		}
	}
	return writeFile(d.filePath, dailyPrgms)
}

func (d EvenProgramRepository) FindAll(ctx context.Context) ([]program.Even, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return nil, err
	}
	return buildEvenPrograms(dailyPgrms), nil
}

func buildEvenPrograms(pr programMap) []program.Even {
	dailies := make([]program.Even, 0, len(pr))
	for hour, pg := range pr {
		dailies = append(dailies, buildEven(pg, hour))
	}
	return dailies
}

func (d EvenProgramRepository) FindByHour(ctx context.Context, hour program.Hour) (program.Even, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return program.Even{}, err
	}
	pgr, ok := dailyPgrms[hour.String()]
	if !ok {
		return program.Even{}, vo.NotFoundError{}
	}
	return buildEven(pgr, hour.String()), nil
}

func buildEven(pgr programData, hour string) program.Even {
	var prog program.Even
	sec, _ := program.ParseSeconds(pgr.Seconds)
	ho, _ := program.ParseHour(hour)
	prog.Hydrate(sec, ho, pgr.Zones)
	return prog
}

func NewEvenProgramRepository(filePath string) EvenProgramRepository {
	return EvenProgramRepository{filePath: filePath}
}
