package http

import (
	"encoding/json"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func FindWeather(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := qh.Handle(r.Context(), app.FindWeatherQuery{})
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		weath, _ := result.(weather.Weather)
		resp := WeatherResponseJson{
			Humidity:    float64(weath.Humidity()),
			IsRaining:   weath.IsRaining(),
			Temperature: float64(weath.Temp()),
		}
		data, _ := json.Marshal(resp)
		WriteResponse(w, http.StatusOK, data)
	}
}
