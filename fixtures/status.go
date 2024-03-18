package fixtures

import (
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

type StatusBuilder struct {
	SystemStartedAt *vo.Time
	Weather         *weather.Weather
	UpdatedAt       *vo.Time
	Active          bool
}

func (b StatusBuilder) Build() status.Status {
	var st status.Status
	start := vo.TimeNow()
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
	st.Hydrate(start, weath, updateAt, b.Active)
	return st
}
