package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/disk/file"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/gpio/relay"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/log/logger"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/mysql"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/telegram"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecutionInTime(t *testing.T) {
	st := status.New()
	log := logger.New()
	log.EnableDebug()
	notifSender := telegram.NewInMemorySender(log)
	executionRepo := file.NewExecutionRepository("./assets/executions.yml")
	zoneRepo := file.NewZoneRepository("./assets/zones.yml")
	relayMang := relay.NewManager()
	conf := mysql.NewConfig("0.0.0.0", "3306", "raspberry", "raspberry", "raspberryWaterSystem")
	myqlRepo := mysql.NewRepository(conf)
	logRepo := mysql.NewExecutionLogRepository(myqlRepo)
	executor := execution.NewExecutor(zoneRepo, relayMang, logRepo, notifSender)

	exec := execution.NewExecutorInTime(executionRepo, executor, st, notifSender)
	err := exec.Execute(time.Now())
	assert.NoError(t, err)
}
