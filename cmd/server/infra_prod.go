//go:build prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/gpio"
)

func temperatureRepository() app.TemperatureRepository {
	repo, _ := gpio.NewBme280TemperatureRepository()
	return repo
}

func pinsExecutor() app.PinExecutor {
	return gpio.NewPinsExecutor()
}
