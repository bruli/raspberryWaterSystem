package telegram

import logger2 "github.com/bruli/raspberryWaterSystem/internal/logger"

type InMemorySender struct {
	log logger2.Logger
}

func NewInMemorySender(log logger2.Logger) *InMemorySender {
	return &InMemorySender{log: log}
}

func (s *InMemorySender) Send(message string) error {
	s.log.Debugf("message sent: %s", message)
	return nil
}
