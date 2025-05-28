package telegram

import (
	"context"
	"errors"
)

type runnerCommand interface {
	CommandName() CommandName
}

type RunnerFunc func(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error

type runner interface {
	Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error
}

type runnerBus struct {
	runners map[CommandName]runner
}

func (r *runnerBus) subscribe(command CommandName, run runner) {
	r.runners[command] = run
}

func (r *runnerBus) handle(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	run, ok := r.runners[cmd.CommandName()]
	if !ok {
		return errors.New("command not found")
	}
	return run.Run(ctx, chatID, msgs, cmd)
}

func newRunnerBus() *runnerBus {
	return &runnerBus{runners: make(map[CommandName]runner)}
}
