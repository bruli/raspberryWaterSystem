package file

import (
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestExecutionRepository_Save(t *testing.T) {
	exec := execution.NewExecutionStub()
	tests := map[string]struct {
		execution         execution.Execution
		err, formattedErr error
	}{
		"it should return error when writer returns error": {
			execution:    exec,
			err:          errors.New("error"),
			formattedErr: fmt.Errorf("failed to save execution: %w", errors.New("error"))},
		"it should save execution": {execution: exec},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			writer := writerMock{}
			repo := &ExecutionRepository{
				repository: &repository{writer: &writer},
			}
			writer.writeFunc = func(d []byte) error {
				return tt.err
			}
			err := repo.Save(tt.execution)
			assert.Equal(t, tt.formattedErr, err)
		})
	}
}

func TestExecutionRepository_GetExecutions(t *testing.T) {
	exec := execution.NewExecutionStub()
	data, _ := yaml.Marshal(&exec)
	tests := map[string]struct {
		err, expectedErr error
		data             []byte
		executions       *execution.Execution
	}{
		"it should return nil when repository returns error": {
			err:        errors.New("error"),
			executions: &execution.Execution{},
		},
		"it should error when unmarshal returns error": {
			expectedErr: errors.New("failed to unmarshal execution"),
			data:        []byte("invalid"),
		},
		"it should return execution": {
			data:       data,
			executions: &exec,
		},
		"it should return empty data": {
			data:       []byte{},
			executions: &execution.Execution{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			reader := readerMock{}
			repo := ExecutionRepository{
				repository: &repository{reader: &reader},
			}
			reader.readFunc = func() ([]byte, error) {
				return tt.data, tt.err
			}

			exec, err := repo.GetExecutions()
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.executions, exec)
		})
	}
}
