package weather

type Weather struct {
	temp, humidity float32
	raining        bool
}

func (w Weather) Temp() float32 {
	return w.temp
}

func (w Weather) Humidity() float32 {
	return w.humidity
}

func (w Weather) IsRaining() bool {
	return w.raining
}

func New(temp, humidity float32, raining bool) Weather {
	return Weather{
		temp:     temp,
		humidity: humidity,
		raining:  raining,
	}
}
