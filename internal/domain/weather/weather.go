package weather

type Temperature float32

func (t Temperature) Float32() float32 {
	return float32(t)
}

type Humidity float32

func (h Humidity) Float32() float32 {
	return float32(h)
}

type Weather struct {
	temp     Temperature
	humidity Humidity
	raining  bool
}

func (w Weather) Temperature() Temperature {
	return w.temp
}

func (w Weather) Humidity() Humidity {
	return w.humidity
}

func (w Weather) IsRaining() bool {
	return w.raining
}

func New(temp Temperature, humidity Humidity, raining bool) Weather {
	return Weather{
		temp:     temp,
		humidity: humidity,
		raining:  raining,
	}
}
