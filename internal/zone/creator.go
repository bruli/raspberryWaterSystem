package zone

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
)

type Creator struct {
	repository      Repository
	relayRepository RelayRepository
	logger          logger.Logger
}

func NewCreator(repository Repository, relayRepository RelayRepository, logger logger.Logger) *Creator {
	return &Creator{repository: repository, relayRepository: relayRepository, logger: logger}
}

func (c *Creator) Create(id, name string, relays []string) error {
	err := c.validateRelays(relays)
	if err != nil {
		c.logger.Fatal(err)
		return err
	}

	zones := c.repository.GetZones()
	current := c.repository.Find(id)
	if current == nil {
		zone, err := New(id, name, relays)
		if err != nil {
			return err
		}
		zones.Add(*zone)
	} else {
		current.update(name, relays)
	}
	if err := c.repository.Save(*zones); err != nil {
		e := fmt.Errorf("failed saving zones: %w", err)
		c.logger.Fatal(e)

		return e
	}

	return nil
}

func (c *Creator) validateRelays(relays []string) error {
	for _, j := range relays {

		if !c.checkRelay(j) {
			return NewInvalidRelay(j)
		}
	}
	return nil
}

func (c *Creator) checkRelay(p string) bool {
	for _, j := range c.relayRepository.Get() {
		if j == p {
			return true
		}
	}

	return false
}
