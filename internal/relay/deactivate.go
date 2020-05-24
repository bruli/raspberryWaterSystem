package relay

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/gpio/relay"
)

type Deactivate struct {
	manag Manager
}

func NewDeactivate(manag Manager) *Deactivate {
	return &Deactivate{manag: manag}
}

func (d *Deactivate) Deactivate() error {
	pins := relay.GetPins()
	err := d.manag.DeactivatePins(pins)
	if err != nil {
		return fmt.Errorf("failed deactivating relays: %w", err)
	}

	return nil
}
