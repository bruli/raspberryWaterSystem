package status

import (
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

type Status struct {
	systemStartedAt vo.Time
	weather         weather.Weather
	updatedAt       *vo.Time
}

func (s Status) SystemStartedAt() vo.Time {
	return s.systemStartedAt
}

func (s Status) Weather() weather.Weather {
	return s.weather
}

func (s Status) UpdatedAt() *vo.Time {
	return s.updatedAt
}

func New(systemStartedAt vo.Time, weather weather.Weather) Status {
	return Status{systemStartedAt: systemStartedAt, weather: weather}
}

func (s *Status) Update(w weather.Weather) {
	s.weather = w
	now := vo.TimeNow()
	s.updatedAt = &now
}

func (s *Status) Hydrate(start vo.Time, we weather.Weather, updated *vo.Time) {
	s.systemStartedAt = start
	s.weather = we
	s.updatedAt = updated
}
