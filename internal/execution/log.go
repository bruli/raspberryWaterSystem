package execution

import (
	"time"
)

type Logs []*Log

func (ls *Logs) Add(l *Log) {
	*ls = append(*ls, l)
}

type Log struct {
	Seconds   uint8
	Zone      string
	CreatedAt time.Time
}

func NewLog(seconds uint8, zone string, createdAt time.Time) *Log {
	return &Log{Seconds: seconds, Zone: zone, CreatedAt: createdAt}
}
