package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getTemperature_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	repo := weather.RepositoryMock{}
	log := logger.LoggerMock{}
	st := status.New()
	getter := weather.NewGetter(&repo)
	router.temperature = newGetTemperature(getter, &log, st)
	server := router.buildServer(config.AuthToken)
	tests := map[string]struct {
		temp, hum float32
		err       error
		code      int
	}{
		"it should return internal server error when repository returns error": {
			err:  errors.New("error"),
			code: http.StatusInternalServerError,
		},
		"it should return temperature data": {
			temp: 25,
			hum:  40,
			code: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/temperature", nil)
			req.Header.Add("Authorization", config.AuthToken)
			assert.NoError(t, err)

			repo.ReadFunc = func() (float32, float32, error) {
				return tt.temp, tt.hum, tt.err
			}
			log.FatalfFunc = func(format string, v ...interface{}) {
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			assert.Equal(t, tt.code, writer.Code)

			body := temperatureBody{}
			err = jsoniter.Unmarshal(writer.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tt.temp, body.Temperature)
			assert.Equal(t, tt.hum, body.Humidity)
		})
	}
}
