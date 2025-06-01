package zone

import (
	"errors"

	"github.com/google/uuid"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const ExecutionSecondsLimit = 300

var (
	ErrInvalidZoneID               = errors.New("zone id can not be empty")
	ErrInvalidZoneName             = errors.New("zone name can not be empty")
	ErrInvalidZoneRelays           = errors.New("zone relays can not be empty")
	ErrInvalidSecondsExecutionZone = errors.New("execution zone has limit 300")
)

const (
	RainingIgnoredReason     = "It's raining!!"
	DeactivatedIgnoredReason = "Deactivated!!"
)

type Zone struct {
	cqs.BasicAggregateRoot
	id, name string
	relays   []Relay
}

func (z *Zone) Id() string {
	return z.id
}

func (z *Zone) Name() string {
	return z.name
}

func (z *Zone) Relays() []Relay {
	return z.relays
}

func New(id, name string, relays []Relay) (*Zone, error) {
	z := Zone{
		BasicAggregateRoot: cqs.NewBasicAggregateRoot(),
		id:                 id,
		name:               name,
		relays:             relays,
	}
	if err := z.validate(); err != nil {
		return nil, err
	}
	return &z, nil
}

func (z *Zone) Hydrate(id, name string, relays []Relay) {
	z.id = id
	z.name = name
	z.relays = relays
}

func (z *Zone) validate() error {
	if len(z.id) == 0 {
		return ErrInvalidZoneID
	}
	if len(z.name) == 0 {
		return ErrInvalidZoneName
	}
	if len(z.relays) == 0 {
		return ErrInvalidZoneRelays
	}
	return nil
}

func (z *Zone) Execute(seconds uint) error {
	if seconds > ExecutionSecondsLimit {
		return ErrInvalidSecondsExecutionZone
	}
	pins := make([]string, len(z.relays))
	for i, p := range z.relays {
		pins[i] = p.pin
	}
	z.Record(Executed{
		BasicEvent: cqs.NewBasicEvent(ExecutedEventName, uuid.New(), z.id),
		ZoneID:     z.id,
		ZoneName:   z.name,
		Seconds:    seconds,
		RelayPins:  pins,
	})
	return nil
}

func (z *Zone) ExecuteWithStatus(active, raining bool, seconds uint) error {
	event := Ignored{
		BasicEvent: cqs.NewBasicEvent(IgnoredEventName, uuid.New(), z.id),
		ZoneName:   z.name,
	}
	switch {
	case !active:
		event.Reason = DeactivatedIgnoredReason
		z.Record(event)
		return nil
	case raining:
		event.Reason = RainingIgnoredReason
		z.Record(event)
		return nil
	default:
		return z.Execute(seconds)
	}
}
