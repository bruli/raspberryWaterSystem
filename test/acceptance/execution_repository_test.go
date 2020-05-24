package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/disk/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutionRepository_save(t *testing.T) {
	config := getConfig()
	exec := execution.NewExecutionStub()
	repo := file.NewExecutionRepository(config.ExecutionsFile)
	err := repo.Save(exec)

	assert.NoError(t, err)

	execut, err := repo.GetExecutions()
	assert.NoError(t, err)
	assert.Equal(t, len(*exec.Daily), len(*execut.Daily))
	assert.Equal(t, len(*exec.Weekly), len(*execut.Weekly))
	assert.Equal(t, len(*exec.Odd), len(*execut.Odd))
	assert.Equal(t, len(*exec.Even), len(*execut.Weekly))
}

func TestExecutionRepository_getExecutions(t *testing.T) {
	config := getConfig()
	repo := file.NewExecutionRepository(config.ExecutionsFile)

	execut, err := repo.GetExecutions()
	assert.NoError(t, err)
	assert.NotNil(t, execut)
}
