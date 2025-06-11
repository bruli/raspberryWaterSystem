//go:build prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/gpio"
)

func temperatureRepository() app.TemperatureRepository {
	repo, _ := gpio.NewBme280TemperatureSensor()
	return repo
}

func pinsExecutor() app.PinExecutor {
	return gpio.NewPinsExecutor()
}

func rainRepository() app.RainRepository {
	return gpio.NewRainSensor()
}
