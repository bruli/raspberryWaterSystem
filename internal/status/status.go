package status

import "time"

type RainStatus struct {
	isRain bool
	value  uint16
}

func (r *RainStatus) Value() uint16 {
	return r.value
}

func (r *RainStatus) IsRain() bool {
	return r.isRain
}

func newRainStatus(isRain bool, value uint16) *RainStatus {
	return &RainStatus{isRain: isRain, value: value}
}

type Status struct {
	systemStarted time.Time
	temperature   float32
	humidity      float32
	onWater       bool
	rain          *RainStatus
}

func (st *Status) OnWater() bool {
	return st.onWater
}

func (st *Status) Humidity() float32 {
	return st.humidity
}

func (st *Status) Temperature() float32 {
	return st.temperature
}

func (st *Status) SystemStarted() time.Time {
	return st.systemStarted
}

func (st *Status) SetTemperature(t float32) {
	st.temperature = t
}

func (st *Status) SetHumidity(h float32) {
	st.humidity = h
}

func (st *Status) SetRain(isRain bool, value uint16) {
	st.rain.isRain = isRain
	st.rain.value = value
}

func (st *Status) Rain() *RainStatus {
	return st.rain
}

func New() *Status {
	return &Status{systemStarted: time.Now(),
		onWater: false,
		rain:    newRainStatus(false, 1023)}
}
