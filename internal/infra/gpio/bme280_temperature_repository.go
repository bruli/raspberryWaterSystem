package gpio

import (
	"context"
	"fmt"
	"sync"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

type Bme280TemperatureRepository struct {
	bus i2c.BusCloser
	sync.Mutex
}

func (b *Bme280TemperatureRepository) Find(ctx context.Context) (weather.Temperature, weather.Humidity, error) {
	select {
	case <-ctx.Done():
		_ = b.bus.Close()
		return 0, 0, ctx.Err()
	default:
		b.Lock()
		sensor, err := bmxx80.NewI2C(b.bus, 0x76, &bmxx80.DefaultOpts)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to create temperature sensor: %w", err)
		}
		defer func() {
			b.Unlock()
			_ = sensor.Halt()
		}()
		var env physic.Env
		if err = sensor.Sense(&env); err != nil {
			return 0, 0, fmt.Errorf("failed to sense temperature: %w", err)
		}

		temp := weather.Temperature(env.Temperature.Celsius())
		hum := weather.Humidity(env.Humidity) / 100000
		return temp, hum, nil
	}
}

func NewBme280TemperatureRepository() (*Bme280TemperatureRepository, error) {
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize temperature sensor host: %v", err)
	}
	bus, err := i2creg.Open("")
	if err != nil {
		return nil, fmt.Errorf("failed to open i2c bus: %v", err)
	}
	return &Bme280TemperatureRepository{bus: bus}, nil
}
