package gpio

import (
	"context"
	"fmt"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
	"sync"
)

type Bme280TemperatureRepository struct {
	bus i2c.BusCloser
	sync.Mutex
}

func (b *Bme280TemperatureRepository) Find(ctx context.Context) (temp, hum float32, err error) {
	select {
	case <-ctx.Done():
		_ = b.bus.Close()
		return 0, 0, ctx.Err()
	default:
		b.Lock()
		sensor, errNew := bmxx80.NewI2C(b.bus, 0x76, &bmxx80.DefaultOpts)
		if errNew != nil {
			err = fmt.Errorf("failed to create sensor: %w", errNew)
			return
		}
		defer func() {
			b.Unlock()
			_ = sensor.Halt()
		}()
		var env physic.Env
		if errSense := sensor.Sense(&env); errSense != nil {
			err = fmt.Errorf("failed to sense temperature: %w", errSense)
			return
		}

		temp = float32(env.Temperature.Celsius())
		hum = float32(env.Humidity) / 100000
		return
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
