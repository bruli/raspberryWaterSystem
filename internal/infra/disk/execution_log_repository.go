package disk

import (
	"context"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type log struct {
	Seconds    int       `json:"seconds"`
	ZoneName   string    `json:"zone_name"`
	ExecutedAt time.Time `json:"executed_at"`
}

type ExecutionLogRepository struct {
	path string
}

func (e ExecutionLogRepository) Save(ctx context.Context, execLogs []program.ExecutionLog) error {
	logs := buildLogs(execLogs)
	return writeJsonFile(e.path, logs)
}

func buildLogs(execLogs []program.ExecutionLog) []log {
	logs := make([]log, len(execLogs))
	for i, l := range execLogs {
		logs[i] = log{
			Seconds:    l.Seconds().Int(),
			ZoneName:   l.ZoneName(),
			ExecutedAt: time.Time(l.ExecutedAt()),
		}
	}
	return logs
}

func (e ExecutionLogRepository) FindAll(ctx context.Context) ([]program.ExecutionLog, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var logs []log
		if err := readJsonFile(e.path, &logs); err != nil {
			return []program.ExecutionLog{}, nil
		}
		return buildExecutionLogs(logs), nil
	}
}

func buildExecutionLogs(logs []log) []program.ExecutionLog {
	execLogs := make([]program.ExecutionLog, len(logs))
	for i, l := range logs {
		var el program.ExecutionLog
		sec, _ := program.ParseSeconds(l.Seconds)
		el.Hydrate(sec, l.ZoneName, l.ExecutedAt)
		execLogs[i] = el
	}
	return execLogs
}

func NewExecutionLogRepository(path string) ExecutionLogRepository {
	return ExecutionLogRepository{path: path}
}
