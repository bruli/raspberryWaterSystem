package status

import (
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

type Status struct {
	systemStartedAt vo.Time
	weather         weather.Weather
	updatedAt       *vo.Time
	active          bool
}

func (s *Status) IsActive() bool {
	return s.active
}

func (s *Status) SystemStartedAt() vo.Time {
	return s.systemStartedAt
}

func (s *Status) Weather() weather.Weather {
	return s.weather
}

func (s *Status) UpdatedAt() *vo.Time {
	return s.updatedAt
}

func (s *Status) Update(w weather.Weather) {
	s.weather = w
	now := vo.TimeNow()
	s.updatedAt = &now
}

func (s *Status) Hydrate(start vo.Time, we weather.Weather, updated *vo.Time, active bool) {
	s.systemStartedAt = start
	s.weather = we
	s.updatedAt = updated
	s.active = active
}

func (s *Status) Activate() {
	s.active = true
}

func (s *Status) Deactivate() {
	s.active = false
}

func New(systemStartedAt vo.Time, weather weather.Weather) Status {
	return Status{systemStartedAt: systemStartedAt, weather: weather, active: true}
}
