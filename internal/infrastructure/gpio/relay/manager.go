package relay

type Manager struct {
	gpio   *gioPins
	relays *relays
}

func (r *Manager) ActivatePins(pins []string) error {
	for _, p := range pins {
		pin := r.relays.getPin(p)
		gpio := r.gpio.getPin(pin)
		gpio.Output().Low()
	}

	return nil
}

func (r *Manager) DeactivatePins(pins []string) error {
	for _, p := range pins {
		pin := r.relays.getPin(p)
		gpio := r.gpio.getPin(pin)
		gpio.Output().High()
	}

	return nil
}

func NewManager() *Manager {
	return &Manager{gpio: newGioPins(), relays: getRelays()}
}
