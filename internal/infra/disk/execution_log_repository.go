package disk

import (
	"context"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Log struct {
	Seconds    int       `json:"seconds"`
	ZoneName   string    `json:"zone_name"`
	ExecutedAt time.Time `json:"executed_at"`
}

type ExecutionLogRepository struct {
	path   string
	tracer trace.Tracer
}

func (e ExecutionLogRepository) Save(ctx context.Context, execLogs []program.ExecutionLog) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := e.tracer.Start(ctx, "ExecutionLogRepository.Save")
		defer span.End()
		logs := buildLogs(execLogs)
		if err := writeJsonFile(e.path, logs); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetStatus(codes.Ok, "execution logs saved")
		return nil
	}
}

func buildLogs(execLogs []program.ExecutionLog) []Log {
	logs := make([]Log, len(execLogs))
	for i, l := range execLogs {
		logs[i] = Log{
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
		_, span := e.tracer.Start(ctx, "ExecutionLogRepository.FindAll")
		defer span.End()
		var logs []Log
		if err := readJsonFile(e.path, &logs); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return []program.ExecutionLog{}, nil
		}
		span.SetStatus(codes.Ok, "execution logs found")
		return buildExecutionLogs(logs), nil
	}
}

func buildExecutionLogs(logs []Log) []program.ExecutionLog {
	execLogs := make([]program.ExecutionLog, len(logs))
	for i, l := range logs {
		var el program.ExecutionLog
		sec, _ := program.ParseSeconds(l.Seconds)
		el.Hydrate(sec, l.ZoneName, l.ExecutedAt)
		execLogs[i] = el
	}
	return execLogs
}

func NewExecutionLogRepository(path string, tracer trace.Tracer) ExecutionLogRepository {
	return ExecutionLogRepository{path: path, tracer: tracer}
}
