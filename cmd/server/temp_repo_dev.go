//go:build !prod
// +build !prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/fake"
)

func temperatureRepository() app.TemperatureRepository {
	return fake.TemperatureRepository{}
}
