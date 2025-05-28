package telegram

import (
	"context"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"strconv"
	"strings"
)

type zoneCommand struct {
	arguments string
}

func (z zoneCommand) CommandName() CommandName {
	return ZoneCommandName
}

type zoneRunner struct {
	ch cqs.CommandHandler
}

func (z zoneRunner) Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	co, _ := cmd.(zoneCommand)
	arguments := strings.Fields(co.arguments)
	if len(arguments) != 3 {
		return fmt.Errorf("invalid arguments")
	}
	id := arguments[0]
	name := arguments[1]
	relays := z.buildRelaysFromArguments(arguments)
	if _, err := z.ch.Handle(ctx, app.CreateZoneCmd{
		ID:       id,
		ZoneName: name,
		Relays:   relays,
	}); err != nil {
		return fmt.Errorf("failed to create zone: %w", err)
	}
	buildMessage(chatID, msgs, fmt.Sprintf("Zone created: %s", name))
	return nil
}

func (z zoneRunner) buildRelaysFromArguments(arguments []string) []int {
	relaysStr := arguments[2]
	split := strings.Split(relaysStr, ",")
	relays := make([]int, len(split))
	for i, s := range split {
		relays[i], _ = strconv.Atoi(s)
	}
	return relays
}

func newZoneRunner(ch cqs.CommandHandler) *zoneRunner {
	return &zoneRunner{ch: ch}
}
