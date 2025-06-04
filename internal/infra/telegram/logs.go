package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type logCommand struct {
	Number int
}

func (l logCommand) CommandName() CommandName {
	return LogCommandName
}

type logRunner struct {
	qh cqs.QueryHandler
}

func (l logRunner) Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	co, _ := cmd.(logCommand)
	result, err := l.qh.Handle(ctx, app.FindExecutionLogsQuery{Limit: co.Number})
	if err != nil {
		return fmt.Errorf("failed to find logs: %w", err)
	}
	logs, _ := result.([]string)
	if len(logs) == 0 {
		buildMessage(chatID, msgs, "No logs found")
		return nil
	}
	for _, lo := range logs {
		buildMessage(chatID, msgs, lo)
	}
	return nil
}

func newLogRunner(qh cqs.QueryHandler) *logRunner {
	return &logRunner{qh: qh}
}
