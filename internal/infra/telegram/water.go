package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type waterCommand struct {
	arguments string
}

func (w waterCommand) CommandName() CommandName {
	return WaterCommandName
}

type waterRunner struct {
	ch cqs.CommandHandler
}

func (w waterRunner) Run(ctx context.Context, _ int64, _ *Messages, cmd runnerCommand) error {
	co, _ := cmd.(waterCommand)
	arguments := strings.Fields(co.arguments)
	if len(arguments) == 0 {
		return fmt.Errorf("invalid arguments")
	}
	zone := arguments[0]
	seconds, err := strconv.Atoi(arguments[1])
	if err != nil {
		return fmt.Errorf("invalid seconds")
	}
	if _, err = w.ch.Handle(ctx, app.ExecuteZoneCmd{
		Seconds: uint(seconds),
		ZoneID:  zone,
	}); err != nil {
		return fmt.Errorf("failed to execute zone: %w", err)
	}
	return nil
}

func newWaterRunner(ch cqs.CommandHandler) *waterRunner {
	return &waterRunner{ch: ch}
}
