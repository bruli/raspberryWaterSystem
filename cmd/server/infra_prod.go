//go:build prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/gpio"
	"go.opentelemetry.io/otel/trace"
)

func temperatureRepository(tracer trace.Tracer) app.TemperatureRepository {
	repo, _ := gpio.NewBme280TemperatureSensor(tracer)
	return repo
}

func pinsExecutor(tracer trace.Tracer) app.PinExecutor {
	return gpio.NewPinsExecutor(tracer)
}

func rainRepository(tracer trace.Tracer) app.RainRepository {
	return gpio.NewRainSensor(tracer)
}
