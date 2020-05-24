package temperature

import "github.com/bruli/raspberryWaterSystem/internal/logger"

type InMemoryReader struct {
	log logger.Logger
}

func NewInMemoryReader(log logger.Logger) *InMemoryReader {
	return &InMemoryReader{log: log}
}

func (i *InMemoryReader) Read() (temp float32, hum float32, err error) {
	i.log.Debug("temperature read.")
	return 22, 50, nil
}
