package telegram

import (
	"context"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type DeactivateCommand struct{}

func (a DeactivateCommand) CommandName() CommandName {
	return DeactivateCommandName
}

type DeactivateRunner struct {
	ch cqs.CommandHandler
}

func (a DeactivateRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	if _, err := a.ch.Handle(ctx, app.ActivateDeactivateServerCmd{Active: false}); err != nil {
		return fmt.Errorf("failed to deaactivate: %w", err)
	}

	buildMessage(chatID, msgs, "Deactivated!!")
	return nil
}

func NewDeactivateRunner(ch cqs.CommandHandler) *DeactivateRunner {
	return &DeactivateRunner{ch: ch}
}
