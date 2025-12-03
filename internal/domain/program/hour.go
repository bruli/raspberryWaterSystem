package program

import (
	"errors"
	"time"
)

const HourLayout = "15:04"

var ErrInvalidExecutionHour = errors.New("invalid execution hour")

type Hour time.Time

func (h Hour) String() string {
	return time.Time(h).Format(HourLayout)
}

func ParseHour(s string) (Hour, error) {
	h, err := time.Parse(HourLayout, s)
	if err != nil {
		return Hour{}, ErrInvalidExecutionHour
	}
	return Hour(h), nil
}
