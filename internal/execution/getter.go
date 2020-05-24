package execution

import logger2 "github.com/bruli/raspberryWaterSystem/internal/logger"

type Getter struct {
	repository Repository
	logger     logger2.Logger
}

func NewGetter(repository Repository, logger logger2.Logger) *Getter {
	return &Getter{repository: repository, logger: logger}
}

func (r *Getter) Get() (*Execution, error) {
	exec, err := r.repository.GetExecutions()
	if err != nil {
		r.logger.Fatalf("failed to get executions: %w", err)
		return nil, err
	}

	return exec, nil
}
