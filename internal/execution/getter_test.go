package execution_test

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRead_Get(t *testing.T) {
	exec := execution.NewExecutionStub()
	tests := map[string]struct {
		err  error
		exec *execution.Execution
	}{
		"it should return error when repository returns error": {err: errors.New("error")},
		"it should return execution":                           {exec: &exec},
		"it should return empty execution":                     {exec: &execution.Execution{}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := execution.RepositoryMock{}
			logger := logger.LoggerMock{}
			getter := execution.NewGetter(&repo, &logger)

			logger.FatalfFunc = func(format string, v ...interface{}) {
			}
			repo.GetExecutionsFunc = func() (*execution.Execution, error) {
				return tt.exec, tt.err
			}

			exec, err := getter.Get()
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.exec, exec)
		})
	}
}
