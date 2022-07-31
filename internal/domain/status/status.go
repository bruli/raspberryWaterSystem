package status

import (
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

type Status struct {
	systemStartedAt time.Time
	weather         weather.Weather
	updatedAt       *time.Time
}

func (s Status) SystemStartedAt() time.Time {
	return s.systemStartedAt
}

func (s Status) Weather() weather.Weather {
	return s.weather
}

func (s Status) UpdatedAt() *time.Time {
	return s.updatedAt
}

func New(systemStartedAt time.Time, weather weather.Weather) Status {
	return Status{systemStartedAt: systemStartedAt, weather: weather}
}

func (s *Status) Update(w weather.Weather) {
	s.weather = w
	now := time.Now()
	s.updatedAt = &now
}

func (s *Status) Hydrate(start time.Time, we weather.Weather, updated *time.Time) {
	s.systemStartedAt = start
	s.weather = we
	s.updatedAt = updated
}
