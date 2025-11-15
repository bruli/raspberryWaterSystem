package fixtures

import "github.com/bruli/raspberryWaterSystem/internal/domain/weather"

type WeatherBuilder struct {
	Temp     *weather.Temperature
	Humidity *weather.Humidity
	Raining  bool
}

func (b WeatherBuilder) Build() weather.Weather {
	temp := weather.Temperature(12)
	if b.Temp != nil {
		temp = *b.Temp
	}
	hum := weather.Humidity(40)
	if b.Humidity != nil {
		hum = *b.Humidity
	}
	return weather.New(temp, hum, b.Raining)
}
