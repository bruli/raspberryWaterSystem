package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type getTemperature struct {
	getter   *weather.Getter
	response *response
	st       *status.Status
}

type temperatureBody struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

func newTemperatureBody(temperature, humidity float32) *temperatureBody {
	return &temperatureBody{Temperature: temperature, Humidity: humidity}
}

func newGetTemperature(getter *weather.Getter, log logger.Logger, st *status.Status) *getTemperature {
	return &getTemperature{getter: getter, response: newResponse(log), st: st}
}

func (t *getTemperature) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	temp, hum, err := t.getter.Get()
	if err != nil {
		t.response.generateJSONErrorResponse(w, err)
		return
	}

	body := newTemperatureBody(temp, hum)
	data, err := jsoniter.Marshal(&body)
	if err != nil {
		t.response.generateJSONErrorResponse(w, err)
		return
	}
	t.st.SetTemperature(temp)
	t.st.SetHumidity(hum)
	t.response.writeJSONResponse(w, http.StatusOK, data)
}
