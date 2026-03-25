package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type MessagePublisher struct {
	token  string
	chatID int
	tracer trace.Tracer
}

func NewMessagePublisher(token string, chatID int, tracer trace.Tracer) *MessagePublisher {
	return &MessagePublisher{token: token, chatID: chatID, tracer: tracer}
}

func (e MessagePublisher) Publish(ctx context.Context, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := e.tracer.Start(ctx, "MessagePublisher.Publish")
		defer span.End()
		bot, err := tgbotapi.NewBotAPI(e.token)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		msg := tgbotapi.NewMessage(int64(e.chatID), message)
		_, err = bot.Send(msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetStatus(codes.Ok, "message published")
		return nil
	}
}
