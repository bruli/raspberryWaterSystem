package ws_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/pkg/ws"
	"github.com/stretchr/testify/require"
)

func TestGetWeather(t *testing.T) {
	weaherResp := http2.WeatherResponseJson{
		Humidity:    10,
		IsRaining:   false,
		Temperature: 20,
	}
	weather := ws.Weather{
		Humidity:    10,
		Temperature: 20,
		IsRaining:   false,
	}
	tests := []struct {
		name                string
		cliErr, expectedErr error
		response            *http.Response
		expectedWeather     ws.Weather
	}{
		{
			name:        "and http client returns an error, then it returns a server error",
			cliErr:      errors.New(""),
			expectedErr: ws.ErrServer,
		},
		{
			name:        "and http client returns an internal server error response, then it returns a remote server error",
			response:    &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:        "and http client returns an ok with invalid response, then it returns a failed to read response error",
			response:    &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, invalidResponse{})},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:            "and http client returns an ok response, then it returns a valid weather object",
			response:        &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, weaherResp)},
			expectedWeather: weather,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a Weather method,
		when it called `+tt.name, func(t *testing.T) {
			t.Parallel()
			cl := &HTTPClientMock{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tt.response, tt.cliErr
				},
			}
			pkg := ws.New(url.URL{}, cl, "token")
			we, err := pkg.GetWeather(context.Background())
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedWeather, we)
		})
	}
}
