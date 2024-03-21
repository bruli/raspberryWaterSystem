package app

import (
	"context"
	"errors"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const FindExecutionLogsQueryName = "findExecutionLogs"

var ErrInvalidExecutionsLogLimit = errors.New("invalid executions log limit")

type FindExecutionLogsQuery struct {
	Limit int
}

func (f FindExecutionLogsQuery) Name() string {
	return FindExecutionLogsQueryName
}

type FindExecutionLogs struct {
	elr ExecutionLogRepository
}

func (f FindExecutionLogs) Handle(ctx context.Context, query cqs.Query) (any, error) {
	q, _ := query.(FindExecutionLogsQuery)
	if q.Limit > maxExecutionLogs {
		return nil, ErrInvalidExecutionsLogLimit
	}
	logs, err := f.elr.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	if q.Limit >= len(logs) {
		return logs, nil
	}
	return logs[len(logs)-q.Limit:], nil
}

func NewFindExecutionLogs(elr ExecutionLogRepository) FindExecutionLogs {
	return FindExecutionLogs{elr: elr}
}
