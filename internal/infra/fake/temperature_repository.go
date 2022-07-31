package fake

import "context"

type TemperatureRepository struct{}

func (t TemperatureRepository) Find(ctx context.Context) (temp, hum float32, err error) {
	return 20.05, 40.3, nil
}
