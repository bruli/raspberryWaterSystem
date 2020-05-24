package file

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"gopkg.in/yaml.v2"
)

type ExecutionRepository struct {
	repository *repository
}

func (ex *ExecutionRepository) Save(e execution.Execution) error {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}
	if err := ex.repository.writer.write(data); err != nil {
		return fmt.Errorf("failed to save execution: %w", err)
	}
	return nil
}

func (ex *ExecutionRepository) GetExecutions() (*execution.Execution, error) {
	exec := execution.Execution{}
	data, err := ex.repository.reader.read()
	if err != nil {
		return &exec, nil
	}

	if err := yaml.Unmarshal(data, &exec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal execution")
	}

	return &exec, nil
}

func NewExecutionRepository(file string) *ExecutionRepository {
	return &ExecutionRepository{repository: newRepository(file)}
}
