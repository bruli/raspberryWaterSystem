package execution

import "github.com/bruli/raspberryWaterSystem/internal/logger"

type Creator struct {
	repository Repository
	logger     logger.Logger
}

func NewCreator(repository Repository, logger logger.Logger) *Creator {
	return &Creator{repository: repository, logger: logger}
}

func (c *Creator) Create(exec Execution) error {
	if err := c.repository.Save(exec); err != nil {
		c.logger.Fatalf("failed to save executions: %w", err)
		return err
	}

	return nil
}
