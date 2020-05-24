package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExecutionLogsHandler_ServeHTTP(t *testing.T) {
	config := getConfig()
	lo, err := execution.NewLogsStub()
	assert.NoError(t, err)
	repo := execution.LogRepositoryMock{}
	log := logger.LoggerMock{}
	router := getRouter()
	router.getExecutionLogs = newGetExecutionLogsHandler(execution.NewReadLogs(&repo, &log), &log)
	server := router.buildServer(config.AuthToken)
	tests := map[string]struct {
		code int
		logs *execution.Logs
		err  error
	}{
		"it should return internal server error when repository returns error": {
			code: http.StatusInternalServerError, err: errors.New("error")},
		"it should return logsBody": {
			logs: &lo,
			code: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/executions/logs", nil)
			req.Header.Add("Authorization", config.AuthToken)
			assert.NoError(t, err)

			repo.GetFunc = func() (*execution.Logs, error) {
				return tt.logs, tt.err
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}
			log.FatalfFunc = func(format string, v ...interface{}) {
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			assert.Equal(t, tt.code, writer.Code)

			if http.StatusOK == writer.Code {
				body := logsBody{}
				err = jsoniter.Unmarshal(writer.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.NotNil(t, body[0].Message)
				assert.NotNil(t, body[0].CreatedAt)
				assert.True(t, execution.MaxLogs >= len(body))
			}
		})
	}
}
