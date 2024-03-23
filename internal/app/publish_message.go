package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const PublishMessageCmdName = "publishMessage"

type PublishMessageCmd struct {
	Message string
}

func (p PublishMessageCmd) Name() string {
	return PublishMessageCmdName
}

type PublishMessage struct {
	mp MessagePublisher
}

func (p PublishMessage) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(PublishMessageCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(PublishMessageCmdName, cmd.Name())
	}
	return nil, p.mp.Publish(ctx, co.Message)
}

func NewPublishMessage(elp MessagePublisher) PublishMessage {
	return PublishMessage{mp: elp}
}
