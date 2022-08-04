//go:build !prod
// +build !prod

package main

import (
	"log"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/fake"
)

func temperatureRepository() app.TemperatureRepository {
	return fake.TemperatureRepository{}
}

func pinsExecutor(log *log.Logger) app.PinExecutor {
	return fake.NewPinsExecutor(log)
}
