//go:build prod
// +build prod

package main

import (
	"log"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/gpio"
)

func temperatureRepository() app.TemperatureRepository {
	return gpio.TemperatureRepository{}
}

func pinsExecutor(log *log.Logger) app.PinExecutor {
	return gpio.PinsExecutor{}
}
