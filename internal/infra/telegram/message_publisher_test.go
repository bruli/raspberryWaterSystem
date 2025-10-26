//go:build integration
// +build integration

package telegram_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/config"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/infra/telegram"
	"github.com/stretchr/testify/require"
)

func TestExecutionLogPublisherPublish(t *testing.T) {
	token, err := config.Value(config.TelegramToken)
	if err != nil {
		t.Fatal(err.Error())
	}
	chatIDStr, err := config.Value(config.TelegramChatID)
	if err != nil {
		t.Fatal(err.Error())
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		t.Fatal(err.Error())
	}
	publisher := telegram.NewMessagePublisher(token, chatID)
	err = publisher.Publish(context.Background(), fixtures.ExecutionLogBuilder{}.Build())
	require.NoError(t, err)
}
