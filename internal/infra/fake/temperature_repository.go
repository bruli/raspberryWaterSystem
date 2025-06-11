package fake

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

type TemperatureRepository struct{}

func (t TemperatureRepository) Find(ctx context.Context) (weather.Temperature, weather.Humidity, error) {
	select {
	case <-ctx.Done():
		return 0, 0, ctx.Err()
	default:
		return weather.Temperature(25.05), weather.Humidity(50.3), nil
	}
}
