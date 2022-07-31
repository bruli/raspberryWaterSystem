package app

import (
	"context"
	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
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

func (f FindWeather) Handle(ctx context.Context, query cqs.Query) (cqs.QueryResult, error) {
	temp, hum, err := f.tr.Find(ctx)
	if err != nil {
		return nil, err
	}
	rain, err := f.rr.Find(ctx)
	if err != nil {
		return nil, err
	}
	return weather.New(temp, hum, rain), nil
}

func NewFindWeather(tr TemperatureRepository, rr RainRepository) FindWeather {
	return FindWeather{tr: tr, rr: rr}
}
