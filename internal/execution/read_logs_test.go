package execution_test

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	logger2 "github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadLogs_Read(t *testing.T) {
	lo, err := execution.NewLogsStub()
	assert.NoError(t, err)
	tests := map[string]struct {
		logs *execution.Logs
		err  error
	}{
		"it should return error when repository returns error": {err: errors.New("error")},
		"it should return logs":                                {logs: &lo},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := execution.LogRepositoryMock{}
			logger := logger2.LoggerMock{}
			read := execution.NewReadLogs(&repo, &logger)

			repo.GetFunc = func() (*execution.Logs, error) {
				return tt.logs, tt.err
			}
			logger.FatalfFunc = func(format string, v ...interface{}) {
			}
			lo, err := read.Read()

			assert.Equal(t, tt.err, err)
			if lo != nil {
				assert.True(t, 0 < len(*lo))
				if execution.MaxLogs < len(*tt.logs) {
					assert.Equal(t, execution.MaxLogs, len(*lo))
				}
			}
		})
	}
}
