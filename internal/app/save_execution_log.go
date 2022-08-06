package app

import (
	"context"
	"fmt"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

const SaveExecutionLogCmdName = "saveExecutionLog"

const maxLogs = 20

type SaveExecutionLogCmd struct {
	ZoneName   string
	Seconds    program.Seconds
	ExecutedAt time.Time
}

func (s SaveExecutionLogCmd) Name() string {
	return SaveExecutionLogCmdName
}

type SaveExecutionLog struct {
	elr ExecutionLogRepository
}

func (s SaveExecutionLog) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(SaveExecutionLogCmd)
	logs, err := s.elr.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	readedLogs := logs
	if len(logs) >= maxLogs {
		readedLogs = logs[len(logs)-maxLogs:]
	}
	log, err := program.NewExecutionLog(co.Seconds, co.ZoneName, co.ExecutedAt)
	if err != nil {
		return nil, SaveExecutionLogError{m: err.Error()}
	}
	readedLogs = append(readedLogs, log)
	return nil, s.elr.Save(ctx, readedLogs)
}

func NewSaveExecutionLog(elr ExecutionLogRepository) SaveExecutionLog {
	return SaveExecutionLog{elr: elr}
}

type SaveExecutionLogError struct {
	m string
}

func (s SaveExecutionLogError) Error() string {
	return fmt.Sprintf("failed to save execution log: %s", s.m)
}
