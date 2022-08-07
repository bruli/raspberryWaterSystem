package program

import (
	"errors"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
)

var (
	ErrEmptyZoneName     = errors.New("empty execution log zone name")
	ErrInvalidExecutedAt = errors.New("invalid execution log executed at")
)

type ExecutionLog struct {
	seconds    Seconds
	zoneName   string
	executedAt vo.Time
}

func (e ExecutionLog) Seconds() Seconds {
	return e.seconds
}

func (e ExecutionLog) ZoneName() string {
	return e.zoneName
}

func (e ExecutionLog) ExecutedAt() vo.Time {
	return e.executedAt
}

func (e ExecutionLog) validate() error {
	if _, err := ParseSeconds(e.seconds.Int()); err != nil {
		return err
	}
	if len(e.zoneName) == 0 {
		return ErrEmptyZoneName
	}
	if e.executedAt.IsZero() {
		return ErrInvalidExecutedAt
	}
	return nil
}

func NewExecutionLog(seconds Seconds, zoneName string, executedAt vo.Time) (ExecutionLog, error) {
	exec := ExecutionLog{seconds: seconds, zoneName: zoneName, executedAt: executedAt}
	if err := exec.validate(); err != nil {
		return ExecutionLog{}, err
	}
	return exec, nil
}

func (e *ExecutionLog) Hydrate(seconds Seconds, zoneName string, executedAt vo.Time) {
	e.seconds = seconds
	e.zoneName = zoneName
	e.executedAt = executedAt
}
