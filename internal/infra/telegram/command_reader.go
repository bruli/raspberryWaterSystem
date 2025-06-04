package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type unknownCommand struct{}

func (u unknownCommand) CommandName() CommandName {
	return ""
}

type CommandReader struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	bus     *runnerBus
}

func (r CommandReader) Read(ctx context.Context, logger *zerolog.Logger) {
	logger.Info().Msg("[TELEGRAM SERVICE] starting command reader")
	msgs := NewMessages()
	for update := range r.updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		var cmd runnerCommand
		switch update.Message.Command() {
		case HelpCommandName.String():
			cmd = helCommand{}
		case StatusCommandName.String():
			cmd = statusCommand{}
		case WeatherCommandName.String():
			cmd = weatherCommand{}
		case LogCommandName.String():
			number, err := strconv.Atoi(update.Message.CommandArguments())
			if err != nil {
				number = 2
			}
			cmd = logCommand{Number: number}
		case ActivateCommandName.String():
			cmd = ActivateCommand{}
		case DeactivateCommandName.String():
			cmd = DeactivateCommand{}
		case WaterCommandName.String():
			cmd = waterCommand{arguments: update.Message.CommandArguments()}
		case ZoneCommandName.String():
			cmd = zoneCommand{arguments: update.Message.CommandArguments()}
		case ProgramsCommandName.String():
			cmd = programsCommand{}
		default:
			cmd = unknownCommand{}
		}
		if err := r.bus.handle(ctx, chatID, msgs, cmd); err != nil {
			buildMessage(chatID, msgs, err.Error())
		}
		for _, j := range msgs.GetMessages() {
			if _, err := r.bot.Send(j); err != nil {
				logger.Error().Err(err).Msg("failed to send message to telegram")
			}
		}

		msgs.Clean()

	}
}

func NewCommandReader(telegramToken string, cqh cqs.QueryHandler, ch cqs.CommandHandler) (*CommandReader, error) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot api: %w", err)
	}
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	return &CommandReader{
		bot:     bot,
		updates: bot.GetUpdatesChan(config),
		bus:     buildRunnerBus(cqh, ch),
	}, nil
}

func buildRunnerBus(qh cqs.QueryHandler, ch cqs.CommandHandler) *runnerBus {
	bus := newRunnerBus()
	bus.subscribe(HelpCommandName, newHelpRunner())
	bus.subscribe(StatusCommandName, newStatusRunner(qh))
	bus.subscribe(WeatherCommandName, newWeatherRunner(qh))
	bus.subscribe(LogCommandName, newLogRunner(qh))
	bus.subscribe(ActivateCommandName, NewActivateRunner(ch))
	bus.subscribe(DeactivateCommandName, NewDeactivateRunner(ch))
	bus.subscribe(WaterCommandName, newWaterRunner(ch))
	bus.subscribe(ZoneCommandName, newZoneRunner(ch))
	bus.subscribe(ProgramsCommandName, newProgramsRunner(qh))

	return bus
}
