//go:build prod
// +build prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/gpio"
)

func temperatureRepository() app.TemperatureRepository {
	return gpio.TemperatureRepository{}
}
