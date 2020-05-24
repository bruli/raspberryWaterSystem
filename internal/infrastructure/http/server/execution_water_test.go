package server

import (
	"bytes"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExecutionWater_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	log := logger.LoggerMock{}
	data := make(chan executionData)
	executionWat := NewExecutionWater(data, &log)
	router.executionWater = executionWat
	server := router.buildServer(config.AuthToken)
	tests := map[string]struct {
		body ExecuteWaterBody
		code int
	}{
		"it should return accepted": {
			body: ExecuteWaterBody{Seconds: 2, Zone: "1"},
			code: http.StatusAccepted,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				for {
					select {
					case <-ticker.C:
						assert.Fail(t, "timeout")
					case d := <-data:
						assert.Equal(t, tt.body.Seconds, d.seconds)
						assert.Equal(t, tt.body.Zone, d.zone)
					}
				}
			}()
			data, err := jsoniter.Marshal(tt.body)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/water", bytes.NewBuffer(data))
			assert.NoError(t, err)
			req.Header.Add("Authorization", config.AuthToken)

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			assert.Equal(t, tt.code, writer.Code)
		})
	}
}
