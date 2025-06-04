package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type ActivateCommand struct{}

func (a ActivateCommand) CommandName() CommandName {
	return ActivateCommandName
}

type ActivateRunner struct {
	ch cqs.CommandHandler
}

func (a ActivateRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	if _, err := a.ch.Handle(ctx, app.ActivateDeactivateServerCmd{Active: true}); err != nil {
		return fmt.Errorf("failed to activate: %w", err)
	}

	buildMessage(chatID, msgs, "Activated!!")
	return nil
}

func NewActivateRunner(ch cqs.CommandHandler) *ActivateRunner {
	return &ActivateRunner{ch: ch}
}
