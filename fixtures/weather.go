package fixtures

import "github.com/bruli/raspberryWaterSystem/internal/domain/weather"

type WeatherBuilder struct {
	Temp, Humidity *float32
	Raining        bool
}

func (b WeatherBuilder) Build() weather.Weather {
	temp := float32(12)
	if b.Temp != nil {
		temp = *b.Temp
	}
	hum := float32(40)
	if b.Humidity != nil {
		hum = *b.Humidity
	}
	return weather.New(temp, hum, b.Raining)
}
