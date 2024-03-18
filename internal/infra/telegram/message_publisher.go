package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MessagePublisher struct {
	token  string
	chatID int
}

func NewMessagePublisher(token string, chatID int) *MessagePublisher {
	return &MessagePublisher{token: token, chatID: chatID}
}

func (e MessagePublisher) Publish(ctx context.Context, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		bot, err := tgbotapi.NewBotAPI(e.token)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(int64(e.chatID), message)
		_, err = bot.Send(msg)
		return err
	}
}
