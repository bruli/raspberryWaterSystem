package execution

import "time"

type Weekly struct {
	Weekday    time.Weekday
	Executions *Programs
}

func NewWeeklyByDay(executions *Programs, day time.Weekday) *Weekly {
	return &Weekly{Executions: executions, Weekday: day}
}
