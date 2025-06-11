package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type statusCommand struct{}

func (s statusCommand) CommandName() CommandName {
	return StatusCommandName
}

type statusRunner struct {
	qh cqs.QueryHandler
}

func (s statusRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	result, err := s.qh.Handle(ctx, app.FindStatusQuery{})
	if err != nil {
		return fmt.Errorf("failed to find status: %w", err)
	}
	st, _ := result.(status.Status)
	buildMessage(chatID, msgs, fmt.Sprintf("System started at: %s", st.SystemStartedAt().Date()))
	buildMessage(chatID, msgs, fmt.Sprintf("Current temperature: %v *C", st.Weather().Temperature()))
	buildMessage(chatID, msgs, fmt.Sprintf("Current humidity: %v", st.Weather().Humidity()))
	buildMessage(chatID, msgs, fmt.Sprintf("Is raining: %v", st.Weather().IsRaining()))
	buildMessage(chatID, msgs, fmt.Sprintf("Active: %v", st.IsActive()))
	if st.UpdatedAt() != nil {
		buildMessage(chatID, msgs, fmt.Sprintf("System updated at: %s", st.UpdatedAt().Date()))
	}
	return nil
}

func newStatusRunner(qh cqs.QueryHandler) *statusRunner {
	return &statusRunner{qh: qh}
}
