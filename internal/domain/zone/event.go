package zone

import "github.com/bruli/raspberryRainSensor/pkg/common/cqs"

const ExecutedEventName = "v1.zone.executed"

type Executed struct {
	cqs.BasicEvent
	ZoneID    string
	ZoneName  string
	Seconds   uint
	RelayPins []string
}
