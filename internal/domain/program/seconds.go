package program

import (
	"errors"
	"time"
)

var ErrZeroProgramSeconds = errors.New("program seconds can not be zero")

type Seconds time.Duration

func (s Seconds) Int() int {
	return int(time.Duration(s).Seconds())
}

func ParseSeconds(i int) (Seconds, error) {
	if i < 0 {
		return 0, ErrZeroProgramSeconds
	}
	return Seconds(time.Duration(i) * time.Second), nil
}
