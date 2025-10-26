package gpio

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3/rpi"
)

var relays = map[string]gpio.PinIO{
	"18": rpi.P1_12,
	"17": rpi.P1_11,
	"23": rpi.P1_16,
	"24": rpi.P1_18,
}
