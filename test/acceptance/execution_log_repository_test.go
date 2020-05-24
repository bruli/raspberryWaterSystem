package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecutionLogRepository(t *testing.T) {
	config := getConfig()
	mysqlConfig := mysql.NewConfig(config.MysqlHost,
		config.MysqlPort,
		config.MysqlUser,
		config.MysqlPass,
		config.MysqlDatabase)
	reposit := mysql.NewRepository(mysqlConfig)
	logRepository := mysql.NewExecutionLogRepository(reposit)
	log := execution.NewLog(20, "1", time.Now())
	err := logRepository.Save(*log)
	assert.NoError(t, err)
	lo, err := logRepository.Get()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(*lo))
}
