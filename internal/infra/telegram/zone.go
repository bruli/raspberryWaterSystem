package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type zoneCommand struct {
	arguments string
}

func (z zoneCommand) CommandName() CommandName {
	return ZoneCommandName
}

type zoneRunner struct {
	ch     cqs.CommandHandler
	tracer trace.Tracer
}

func (z zoneRunner) Run(ctx context.Context, chatID int64, msgs *Messages, cmd runnerCommand) error {
	ctx, span := z.tracer.Start(ctx, "zoneRunner.Run")
	defer span.End()
	co, _ := cmd.(zoneCommand)
	arguments := strings.Fields(co.arguments)
	if len(arguments) != 3 {
		err := fmt.Errorf("invalid arguments")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	id := arguments[0]
	name := arguments[1]
	relays := z.buildRelaysFromArguments(arguments)
	if _, err := z.ch.Handle(ctx, app.CreateZoneCmd{
		ID:       id,
		ZoneName: name,
		Relays:   relays,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to create zone: %w", err)
	}
	buildMessage(chatID, msgs, fmt.Sprintf("Zone created: %s", name))
	span.SetStatus(codes.Ok, "zone created")
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

func newZoneRunner(ch cqs.CommandHandler, tracer trace.Tracer) *zoneRunner {
	return &zoneRunner{ch: ch, tracer: tracer}
}
