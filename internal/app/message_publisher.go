package app

import (
	"context"
)

//go:generate go tool moq -out zmock_message_publisher_test.go --pkg app_test . MessagePublisher

type MessagePublisher interface {
	Publish(ctx context.Context, message string) error
}
