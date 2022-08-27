package ws

import "github.com/bruli/raspberryRainSensor/pkg/common/vo"

type Log struct {
	ExecutedAt vo.Time
	Seconds    int
	ZoneName   string
}
