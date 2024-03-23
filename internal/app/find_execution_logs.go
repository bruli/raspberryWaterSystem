package app

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

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
	q, ok := query.(FindExecutionLogsQuery)
	if !ok {
		return nil, cqs.NewInvalidQueryError(FindExecutionLogsQueryName, query.Name())
	}
	if q.Limit > maxExecutionLogs {
		return nil, ErrInvalidExecutionsLogLimit
	}
	logs, err := f.elr.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	if q.Limit >= len(logs) {
		return f.orderLogs(logs), nil
	}
	filtered := logs[len(logs)-q.Limit:]

	return f.orderLogs(filtered), nil
}

func (f FindExecutionLogs) orderLogs(logs []program.ExecutionLog) []program.ExecutionLog {
	sort.Slice(logs, func(i, j int) bool {
		return time.Time(logs[i].ExecutedAt()).After(time.Time(logs[j].ExecutedAt()))
	})
	return logs
}

func NewFindExecutionLogs(elr ExecutionLogRepository) FindExecutionLogs {
	return FindExecutionLogs{elr: elr}
}
