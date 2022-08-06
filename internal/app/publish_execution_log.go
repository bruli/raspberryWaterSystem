package app

import (
	"context"
	"fmt"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

const PublishExecutionLogCmdName = "publishExecutionLog"

type PublishExecutionLogCmd struct {
	ZoneName   string
	Seconds    program.Seconds
	ExecutedAt time.Time
}

func (p PublishExecutionLogCmd) Name() string {
	return PublishExecutionLogCmdName
}

type PublishExecutionLog struct {
	elp ExecutionLogPublisher
}

func (p PublishExecutionLog) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(PublishExecutionLogCmd)
	execLog, err := program.NewExecutionLog(co.Seconds, co.ZoneName, co.ExecutedAt)
	if err != nil {
		return nil, PublishExecutionLogError{m: err.Error()}
	}
	return nil, p.elp.Publish(ctx, execLog)
}

func NewPublishExecutionLog(elp ExecutionLogPublisher) PublishExecutionLog {
	return PublishExecutionLog{elp: elp}
}

type PublishExecutionLogError struct {
	m string
}

func (p PublishExecutionLogError) Error() string {
	return fmt.Sprintf("failed to publish execution log: %s", p.m)
}
