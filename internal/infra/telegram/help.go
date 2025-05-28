package telegram

import (
	"context"
	"fmt"
)

type helCommand struct{}

func (h helCommand) CommandName() CommandName {
	return HelpCommandName
}

type helpRunner struct{}

func (h helpRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		buildMessage(chatID, msgs, "Available commands:")
		for _, c := range initCommands() {
			buildMessage(chatID, msgs, fmt.Sprintf("%s -> %s, %q", c.name.String(), c.syntax, c.description))
		}
	}
	return nil
}

func newHelpRunner() *helpRunner {
	return &helpRunner{}
}
