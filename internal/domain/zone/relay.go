package zone

type Relay struct {
	key, pin string
}

func (r Relay) Key() string {
	return r.key
}

func (r Relay) Pin() string {
	return r.pin
}

func NewRelay(key string, pin string) Relay {
	return Relay{key: key, pin: pin}
}
