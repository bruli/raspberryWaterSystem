package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/jsontime"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getConfig() *Config {
	return NewConfig(
		"http:192.168.1.10",
		"zones.yml",
		"auth_token",
		"executions.yml",
		"locahost",
		"3306",
		"user",
		"pass",
		"database",
		"telegramToken",
		"http:192.168.1.10",
		122345,
		true)
}
func getRouter() *router {
	return newRouter(
		&homePage{},
		&createZone{},
		&getZones{},
		&getExecutionLogs{},
		&createExecution{},
		&getExecutions{},
		&ExecutionWater{},
		&getTemperature{},
		&removeZone{},
		&getRain{})
}

func TestHomePageHandler_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	logger := logger.LoggerMock{}
	status := status.New()
	router.homepage = newHomePage(&logger, status)
	server := router.buildServer(config.AuthToken)
	t.Run("it should return homepage", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Authorization", config.AuthToken)
		assert.NoError(t, err)

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)

		body := homePageResponse{}
		err = jsoniter.Unmarshal(writer.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, status.SystemStarted().Format(jsontime.Layout), body.SystemStarted.ToString())
		assert.Equal(t, status.Humidity(), body.Humidity)
		assert.Equal(t, status.Temperature(), body.Temperature)
		assert.Equal(t, status.OnWater(), body.OnWater)
		assert.Equal(t, status.Rain().Value(), body.Rain.Value)
		assert.Equal(t, status.Rain().IsRain(), body.Rain.IsRaining)
	})
}
