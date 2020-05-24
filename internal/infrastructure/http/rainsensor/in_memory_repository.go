package rainsensor

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
)

type InMemoryRepository struct {
	log logger.Logger
}

func (i *InMemoryRepository) Get() (rain.Rain, error) {
	i.log.Debug("getting rain data")

	return rain.New(false, 1023), nil
}

func NewInMemoryRepository(log logger.Logger) *InMemoryRepository {
	return &InMemoryRepository{log: log}
}
