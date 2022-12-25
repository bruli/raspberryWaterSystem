package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

type ExecutionLogPublisher struct {
	token  string
	chatID int
}

func NewExecutionLogPublisher(token string, chatID int) *ExecutionLogPublisher {
	return &ExecutionLogPublisher{token: token, chatID: chatID}
}

func (e ExecutionLogPublisher) Publish(ctx context.Context, execLog program.ExecutionLog) error {
	bot, err := tgbotapi.NewBotAPI(e.token)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("%s zone executed during %vs", execLog.ZoneName(), execLog.Seconds().Int())
	msg := tgbotapi.NewMessage(int64(e.chatID), message)
	_, err = bot.Send(msg)
	return err
}
