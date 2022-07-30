package zone

import "errors"

const (
	OneRelayID = iota + 1
	TwoRelayID
	ThreeRelayID
	FourRelayID
	FiveRelayID
	SixRelayID
)

type RelayID int

func (i RelayID) Int() int {
	return int(i)
}

var ErrUnknownRelay = errors.New("unknown relay")

var enabledRelays = map[RelayID]string{
	OneRelayID:   "18",
	TwoRelayID:   "24",
	ThreeRelayID: "23",
	FourRelayID:  "25",
	FiveRelayID:  "17",
	SixRelayID:   "27",
}

type Relay struct {
	id  RelayID
	pin string
}

func (r Relay) Id() RelayID {
	return r.id
}

func (r Relay) Pin() string {
	return r.pin
}

func newRelay(id RelayID) (Relay, error) {
	pin, ok := enabledRelays[id]
	if !ok {
		return Relay{}, ErrUnknownRelay
	}
	return Relay{id: id, pin: pin}, nil
}

func ParseRelay(i int) (Relay, error) {
	return newRelay(RelayID(i))
}
