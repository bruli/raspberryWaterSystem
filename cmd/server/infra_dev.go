//go:build !prod

package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/infra/fake"
	"go.opentelemetry.io/otel/trace"
)

func temperatureRepository(_ trace.Tracer) app.TemperatureRepository {
	return fake.TemperatureRepository{}
}

func pinsExecutor(_ trace.Tracer) app.PinExecutor {
	return fake.NewPinsExecutor()
}

func rainRepository(_ trace.Tracer) app.RainRepository {
	return fake.RainRepository{}
}
