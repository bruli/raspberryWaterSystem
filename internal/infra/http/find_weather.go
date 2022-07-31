package http

import (
	"encoding/json"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

func FindWeather(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := qh.Handle(r.Context(), app.FindWeatherQuery{})
		if err != nil {
			httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		weath, _ := result.(weather.Weather)
		resp := WeatherResponseJson{
			Humidity:    float64(weath.Humidity()),
			IsRaining:   weath.IsRaining(),
			Temperature: float64(weath.Temp()),
		}
		data, _ := json.Marshal(resp)
		httpx.WriteResponse(w, http.StatusOK, data)
	}
}
