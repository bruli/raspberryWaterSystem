package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const FindWeatherQueryName = "findWeather"

type FindWeatherQuery struct{}

func (f FindWeatherQuery) Name() string {
	return FindWeatherQueryName
}

type FindWeather struct {
	tr     TemperatureRepository
	rr     RainRepository
	tracer trace.Tracer
}

func (f FindWeather) Handle(ctx context.Context, _ cqs.Query) (any, error) {
	ctx, span := f.tracer.Start(ctx, "FindWeather")
	defer span.End()
	temp, hum, err := f.tr.Find(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	rain, err := f.rr.Find(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetStatus(codes.Ok, "weather found")
	return weather.New(temp, hum, rain), nil
}

func NewFindWeather(tr TemperatureRepository, rr RainRepository, tracer trace.Tracer) FindWeather {
	return FindWeather{tr: tr, rr: rr, tracer: tracer}
}
