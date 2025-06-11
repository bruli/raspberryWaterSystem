package gpio

import (
	"context"
	"sync"

	"github.com/stianeikeland/go-rpio/v4"
)

const RainReference = 600

type RainSensor struct {
	sync.RWMutex
}

func (r *RainSensor) Find(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		r.RLock()
		defer r.RUnlock()
		if err := rpio.Open(); err != nil {
			return false, err
		}

		defer func() {
			_ = rpio.Close()
		}()

		if err := rpio.SpiBegin(rpio.Spi0); err != nil {
			return false, err
		}

		rpio.SpiSpeed(1000000)
		rpio.SpiChipSelect(0)
		channel := byte(0)
		data := []byte{1, (8 + channel) << 4, 0}

		rpio.SpiExchange(data)

		value := int(data[1]&3)<<8 + int(data[2])
		defer rpio.SpiEnd(rpio.Spi0)

		return value > RainReference, nil
	}
}

func NewRainSensor() *RainSensor {
	return &RainSensor{}
}
