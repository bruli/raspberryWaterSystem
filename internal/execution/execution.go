package execution

import "time"

type Execution struct {
	Daily  *Programs
	Weekly *WeeklyPrograms
	Odd    *Programs
	Even   *Programs
	Temp   *TemperaturePrograms
}

func New(daily *Programs, weekly *WeeklyPrograms, odd, even *Programs, temp *TemperaturePrograms) (*Execution, error) {
	if daily == nil && weekly == nil && odd == nil && even == nil && temp == nil {
		return nil, NewInvalidCreateExecution("execution can not be empty")
	}

	if len(daily.GetPrograms()) == 0 &&
		len(weekly.getPrograms()) == 0 &&
		len(odd.GetPrograms()) == 0 &&
		len(even.GetPrograms()) == 0 &&
		len(*temp) == 0 {
		return nil, NewInvalidCreateExecution("any program are into execution")
	}

	return &Execution{Daily: daily, Weekly: weekly, Odd: odd, Even: even, Temp: temp}, nil
}

func (ex *Execution) GetToday(t time.Time, temp float32) *Programs {
	day := t.Day()
	weekday := t.Weekday()
	hour := t.Format("15:04")

	execPrgms := Programs{}

	if ex.Daily != nil {
		dailyPrgms := ex.Daily.GetPrograms()
		if len(dailyPrgms) != 0 {
			daily := dailyPrgms[hour]
			if daily != nil {
				for _, d := range *daily {
					execPrgms.Add(d)
				}
			}
		}
	}

	if ex.Weekly != nil {
		weeklyPrgms := ex.Weekly.getByDay(weekday)
		if weeklyPrgms != nil {
			weekPrgms := weeklyPrgms.GetPrograms()
			weekByHour := weekPrgms[hour]
			if weekByHour != nil {
				for _, p := range *weekByHour {
					execPrgms.Add(p)
				}
			}
		}
	}

	mod := day % 2
	if mod == 0 {
		if ex.Odd != nil {
			oddPrgms := ex.Odd.GetPrograms()
			if len(oddPrgms) != 0 {
				odd := oddPrgms[hour]
				if odd != nil {
					for _, p := range *odd {
						execPrgms.Add(p)
					}
				}
			}
		}
	} else {
		if ex.Even != nil {
			evenPrgms := ex.Even.GetPrograms()
			if len(evenPrgms) != 0 {
				even := evenPrgms[hour]
				if even != nil {
					for _, p := range *even {
						execPrgms.Add(p)
					}
				}
			}
		}
	}

	if ex.Temp != nil {
		tempPrgms := ex.Temp.GetPrograms(temp).GetPrograms()
		if len(tempPrgms) != 0 {
			te := tempPrgms[hour]
			if te != nil {
				for _, p := range *te {
					execPrgms.Add(p)
				}
			}
		}
	}
	return &execPrgms
}
