package ws

import (
	"context"
	"net/http"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

type Weather struct {
	Humidity, Temperature float64
	IsRaining             bool
}

func GetWeather(cl client) WeatherFunc {
	return func(ctx context.Context) (Weather, error) {
		url := cl.serverURL.String() + "/weather"
		resp, err := buildRequestAndSend(ctx, http.MethodGet, nil, url, cl.token, cl.cl)
		if err != nil {
			return Weather{}, ErrServer
		}
		defer func() { _ = resp.Body.Close() }()
		switch resp.StatusCode {
		case http.StatusInternalServerError:
			return Weather{}, ErrRemoteServerErr
		default:
			var schema http2.WeatherResponseJson
			if err = readResponse(resp, &schema); err != nil {
				return Weather{}, ErrFailedToReadResponse
			}
			return Weather{
				Humidity:    schema.Humidity,
				Temperature: schema.Temperature,
				IsRaining:   schema.IsRaining,
			}, nil
		}
	}
}
