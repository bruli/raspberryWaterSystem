package status

import (
	"errors"
	"time"
)

var (
	ErrInvalidSunrise = errors.New("sunrise time is required")
	ErrInvalidSunset  = errors.New("sunset time is required")
)

type Light struct {
	sunrise time.Time
	sunset  time.Time
}

func (l Light) IsDay(t time.Time) bool {
	return l.sunrise.Before(t) && l.sunset.After(t)
}

func (l Light) validate() error {
	switch {
	case l.sunrise.IsZero():
		return ErrInvalidSunrise
	case l.sunset.IsZero():
		return ErrInvalidSunset
	default:
		return nil
	}
}

func NewLight(sunrise time.Time, sunset time.Time) (*Light, error) {
	l := Light{sunrise: sunrise, sunset: sunset}
	if err := l.validate(); err != nil {
		return nil, err
	}
	return &l, nil
}
