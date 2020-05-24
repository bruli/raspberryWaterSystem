package execution

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"time"
)

type ExecutorInTime struct {
	repository  Repository
	exec        *Executor
	st          *status.Status
	notifSender NotificationSender
}

func NewExecutorInTime(repository Repository,
	exec *Executor,
	st *status.Status,
	notifSender NotificationSender) *ExecutorInTime {
	return &ExecutorInTime{repository: repository, exec: exec, st: st, notifSender: notifSender}
}

func (e *ExecutorInTime) Execute(t time.Time) error {
	exec, err := e.repository.GetExecutions()
	if err != nil {
		return fmt.Errorf("failed to get executions: %w", err)
	}
	prgms := exec.GetToday(t)
	for _, prg := range *prgms {
		zons := prg.Executions.Zones
		for _, z := range zons {
			sec := uint8(prg.Seconds.Seconds())
			if e.st.Rain().IsRain() {
				msg := fmt.Sprintf("Is raining. Cannot execute zone %s during %s.", z, prg.Seconds.String())
				if err := e.notifSender.Send(msg); err != nil {
					return fmt.Errorf("failed to send execution notification: %w", err)
				}
				continue
			}
			err := e.exec.Execute(sec, z)
			if err != nil {
				return fmt.Errorf("failed to execute zone '%s': %w", z, err)
			}
		}
	}
	return nil
}
