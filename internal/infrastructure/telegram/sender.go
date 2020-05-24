package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Sender struct {
	token  string
	chatID int64
}

func NewSender(token string, chatID int64) *Sender {
	return &Sender{token: token, chatID: chatID}
}

func (t *Sender) Send(message string) error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err = bot.Send(msg)
	return err
}
