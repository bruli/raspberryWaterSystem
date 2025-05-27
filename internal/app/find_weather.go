package app

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const FindWeatherQueryName = "findWeather"

type FindWeatherQuery struct{}

func (f FindWeatherQuery) Name() string {
	return FindWeatherQueryName
}

type FindWeather struct {
	tr TemperatureRepository
	rr RainRepository
}

func (f FindWeather) Handle(ctx context.Context, _ cqs.Query) (any, error) {
	temp, hum, err := f.tr.Find(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to find temperature")
	}
	rain, err := f.rr.Find(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to find rain")
	}
	return weather.New(temp, hum, rain), nil
}

func NewFindWeather(tr TemperatureRepository, rr RainRepository) FindWeather {
	return FindWeather{tr: tr, rr: rr}
}
