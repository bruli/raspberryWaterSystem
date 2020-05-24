package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExecutions_ServeHTTP(t *testing.T) {
	execut := execution.NewExecutionStub()
	tests := map[string]struct {
		code int
		err  error
		exec *execution.Execution
	}{
		"it should return internal server error when getter returns error": {
			err:  errors.New("error"),
			code: 500,
		},
		"it should return executions": {
			code: 200,
			exec: &execut,
		},
		"it should return empty executions": {
			code: 200,
			exec: &execution.Execution{},
		},
	}
	config := getConfig()
	router := getRouter()
	repo := execution.RepositoryMock{}
	log := logger.LoggerMock{}
	router.getExecutions = newGetExecutions(execution.NewGetter(&repo, &log), &log)
	server := router.buildServer(config.AuthToken)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/executions", nil)
			assert.NoError(t, err)
			req.Header.Add("Authorization", config.AuthToken)

			writer := httptest.NewRecorder()

			log.FatalfFunc = func(format string, v ...interface{}) {
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}
			repo.GetExecutionsFunc = func() (*execution.Execution, error) {
				return tt.exec, tt.err
			}

			server.ServeHTTP(writer, req)

			assert.Equal(t, tt.code, writer.Code)
			if tt.exec != nil {
				body := ExecutionBody{}
				data, err := ioutil.ReadAll(writer.Body)
				assert.NoError(t, err)
				err = jsoniter.Unmarshal(data, &body)
				assert.NoError(t, err)
				if tt.exec.Daily != nil {
					assert.Equal(t, len(*tt.exec.Daily), len(*body.Daily))
				}
			}
		})
	}
}
