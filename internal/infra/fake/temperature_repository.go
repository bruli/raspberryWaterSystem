package fake

import "context"

type TemperatureRepository struct{}

func (t TemperatureRepository) Find(ctx context.Context) (temp, hum float32, err error) {
	select {
	case <-ctx.Done():
		return 0, 0, ctx.Err()
	default:
		return 25.05, 50.3, nil
	}
}
