package disk

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type OddProgramRepository struct {
	filePath string
}

func (d OddProgramRepository) Save(ctx context.Context, programs []program.Odd) error {
	dailyPrgms := make(programMap)
	for _, pr := range programs {
		dailyPrgms[pr.Hour().String()] = programData{
			Seconds: pr.Seconds().Int(),
			Zones:   pr.Zones(),
		}
	}
	return writeFile(d.filePath, dailyPrgms)
}

func (d OddProgramRepository) FindAll(ctx context.Context) ([]program.Odd, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return nil, err
	}
	return buildOddPrograms(dailyPgrms), nil
}

func buildOddPrograms(pr programMap) []program.Odd {
	dailies := make([]program.Odd, 0, len(pr))
	for hour, pg := range pr {
		dailies = append(dailies, buildOdd(pg, hour))
	}
	return dailies
}

func (d OddProgramRepository) FindByHour(ctx context.Context, hour program.Hour) (program.Odd, error) {
	dailyPgrms := make(programMap)
	if err := readFile(d.filePath, &dailyPgrms); err != nil {
		return program.Odd{}, err
	}
	pgr, ok := dailyPgrms[hour.String()]
	if !ok {
		return program.Odd{}, vo.NotFoundError{}
	}
	return buildOdd(pgr, hour.String()), nil
}

func buildOdd(pgr programData, hour string) program.Odd {
	var prog program.Odd
	sec, _ := program.ParseSeconds(pgr.Seconds)
	ho, _ := program.ParseHour(hour)
	prog.Hydrate(sec, ho, pgr.Zones)
	return prog
}

func NewOddProgramRepository(filePath string) OddProgramRepository {
	return OddProgramRepository{filePath: filePath}
}
