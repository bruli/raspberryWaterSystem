package worker

import (
	"context"
	"time"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func ExecutionInTime(ctx context.Context, qh cqs.QueryHandler, ch cqs.CommandHandler, now vo.Time) error {
	st, err := readingCurrentStatus(ctx, qh)
	if err != nil {
		return err
	}
	temp := st.Weather().Temp()
	prgms, err := findingPrograms(ctx, qh, now, temp)
	if err != nil {
		return err
	}
	if err = executeDaily(ctx, prgms, ch, now); err != nil {
		return err
	}
	if err = executeOddEven(ctx, now, prgms, ch); err != nil {
		return err
	}
	if err = executeWeekly(ctx, prgms, ch, now); err != nil {
		return err
	}
	if err = executeTemperature(ctx, prgms, ch, now, temp); err != nil {
		return err
	}
	return nil
}

func findingPrograms(ctx context.Context, qh cqs.QueryHandler, now vo.Time, temp float32) (app.ProgramsInTime, error) {
	resultPrmgs, err := qh.Handle(ctx, app.FindProgramsInTimeQuery{
		On:          now,
		Temperature: temp,
	})
	if err != nil {
		return app.ProgramsInTime{}, FindProgramsError{err: err}
	}
	prgms, _ := resultPrmgs.(app.ProgramsInTime)
	return prgms, nil
}

func readingCurrentStatus(ctx context.Context, qh cqs.QueryHandler) (status.Status, error) {
	resultSt, err := qh.Handle(ctx, app.FindStatusQuery{})
	if err != nil {
		return status.Status{}, ReadCurrentStatusError{err: err}
	}
	st, _ := resultSt.(status.Status)
	return st, nil
}

func executeTemperature(ctx context.Context, prgms app.ProgramsInTime, ch cqs.CommandHandler, now vo.Time, currentTemp float32) error {
	if prgms.Temperature != nil {
		if currentTemp >= prgms.Temperature.Temperature() {
			for _, pr := range prgms.Temperature.Programs() {
				if err := executeProgram(ctx, ch, pr, now); err != nil {
					return ExecuteTemperatureError{err: err}
				}
			}
		}
	}
	return nil
}

func executeWeekly(ctx context.Context, prgms app.ProgramsInTime, ch cqs.CommandHandler, now vo.Time) error {
	if prgms.Weekly != nil {
		if time.Time(now).Weekday().String() == prgms.Weekly.WeekDay().String() {
			for _, pr := range prgms.Weekly.Programs() {
				if err := executeProgram(ctx, ch, pr, now); err != nil {
					return ExecuteWeeklyError{err: err}
				}
			}
		}
	}
	return nil
}

func executeOddEven(ctx context.Context, now vo.Time, prgms app.ProgramsInTime, ch cqs.CommandHandler) error {
	var oddEvenPrgms *program.Program

	switch {
	case isEven(time.Time(now)) && prgms.Even != nil:
		oddEvenPrgms = prgms.Even
	default:
		oddEvenPrgms = prgms.Odd
	}
	if oddEvenPrgms != nil {
		if err := executeProgram(ctx, ch, *oddEvenPrgms, now); err != nil {
			return ExecuteOddEvenError{err: err}
		}
	}
	return nil
}

func executeDaily(ctx context.Context, prgms app.ProgramsInTime, ch cqs.CommandHandler, now vo.Time) error {
	if prgms.Daily != nil {
		if err := executeProgram(ctx, ch, *prgms.Daily, now); err != nil {
			return ExecuteDailyError{err: err}
		}
	}
	return nil
}

func isEven(now time.Time) bool {
	day := now.Day()
	rest := day % 2
	return rest == 0
}

func executeProgram(ctx context.Context, ch cqs.CommandHandler, prg program.Program, now vo.Time) error {
	nowHour := now.HourStr()
	if nowHour == prg.Hour().String() {
		for _, exec := range prg.Executions() {
			for _, zo := range exec.Zones() {
				if _, err := ch.Handle(ctx, app.ExecuteZoneWithStatusCmd{
					Seconds: uint(exec.Seconds().Int()),
					ZoneID:  zo,
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
