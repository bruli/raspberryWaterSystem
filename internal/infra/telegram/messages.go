package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Messages struct {
	msgs []tgbotapi.MessageConfig
}

func (m *Messages) AddMessage(msg tgbotapi.MessageConfig) {
	m.msgs = append(m.msgs, msg)
}

func (m *Messages) Clean() {
	m.msgs = nil
}

func (m *Messages) GetMessages() []tgbotapi.MessageConfig {
	return m.msgs
}

func NewMessages() *Messages {
	return &Messages{msgs: make([]tgbotapi.MessageConfig, 0)}
}

func buildMessage(chatID int64, msgs *Messages, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msgs.AddMessage(msg)
}
