package gpio

import (
	"context"

	"github.com/d2r2/go-dht"
)

type TemperatureRepository struct{}

func (t TemperatureRepository) Find(ctx context.Context) (temp, hum float32, err error) {
	sensorType := dht.DHT11
	pin := 4
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	if err != nil {
		return 0.00, 0.00, err
	}
	return temperature, humidity, nil
}
