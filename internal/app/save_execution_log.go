package app

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

const SaveExecutionLogCmdName = "saveExecutionLog"

const maxExecutionLogs = 20

type SaveExecutionLogCmd struct {
	ZoneName   string
	Seconds    program.Seconds
	ExecutedAt vo.Time
}

func (s SaveExecutionLogCmd) Name() string {
	return SaveExecutionLogCmdName
}

type SaveExecutionLog struct {
	elr ExecutionLogRepository
}

func (s SaveExecutionLog) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(SaveExecutionLogCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(SaveExecutionLogCmdName, cmd.Name())
	}
	logs, err := s.elr.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	readedLogs := logs
	if len(logs) >= maxExecutionLogs {
		readedLogs = logs[len(logs)-maxExecutionLogs:]
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
