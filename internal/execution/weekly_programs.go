package execution

import "time"

type WeeklyPrograms []*Weekly

func (p *WeeklyPrograms) Add(w *Weekly) {
	*p = append(*p, w)
}

func (p *WeeklyPrograms) getPrograms() map[time.Weekday]*Programs {
	e := make(map[time.Weekday]*Programs)
	for _, j := range *p {
		e[j.Weekday] = j.Executions
	}

	return e
}

func (p *WeeklyPrograms) getByDay(weekday time.Weekday) *Programs {
	m := p.getPrograms()
	return m[weekday]
}
