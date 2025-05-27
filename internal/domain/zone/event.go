package zone

import "github.com/bruli/raspberryWaterSystem/pkg/cqs"

const (
	ExecutedEventName = "v1.zone.executed"
	IgnoredEventName  = "v1.zone.ignored"
)

type Executed struct {
	cqs.BasicEvent
	ZoneID    string
	ZoneName  string
	Seconds   uint
	RelayPins []string
}

type Ignored struct {
	cqs.BasicEvent
	ZoneName string
	Reason   string
}
