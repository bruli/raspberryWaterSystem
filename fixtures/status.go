package fixtures

import (
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

type StatusBuilder struct {
	SystemStartedAt *time.Time
	Weather         *weather.Weather
	UpdatedAt       *time.Time
}

func (b StatusBuilder) Build() status.Status {
	var st status.Status
	start := time.Now()
	if b.SystemStartedAt != nil {
		start = *b.SystemStartedAt
	}
	weath := WeatherBuilder{}.Build()
	if b.Weather != nil {
		weath = *b.Weather
	}
	update := start.Add(3 * time.Hour)
	updateAt := &update
	if b.UpdatedAt != nil {
		updateAt = b.UpdatedAt
	}
	st.Hydrate(start, weath, updateAt)
	return st
}
