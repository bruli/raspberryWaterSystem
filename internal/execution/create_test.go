package execution_test

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	logger2 "github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		execution execution.Execution
		err       error
	}{
		"it should return error when repository returns error": {
			execution: execution.NewExecutionStub(),
			err:       errors.New("error"),
		},
		"it should create execution": {
			execution: execution.NewExecutionStub(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := execution.RepositoryMock{}
			logger := logger2.LoggerMock{}
			creat := execution.NewCreator(&repo, &logger)

			repo.SaveFunc = func(e execution.Execution) error {
				return tt.err
			}

			logger.FatalfFunc = func(format string, v ...interface{}) {
			}

			err := creat.Create(tt.execution)

			assert.Equal(t, tt.err, err)
		})
	}
}
