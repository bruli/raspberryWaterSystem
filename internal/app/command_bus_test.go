package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestCommandBusHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name, cmdName  string
		handler        cqs.CommandHandler
		cmd            cqs.Command
		expectedEvents []cqs.Event
		expectedErr    error
	}{
		{
			name:        "with a not subscribed command, then it returns an unsubscribed command error",
			cmdName:     "command",
			handler:     commandHandler{},
			cmd:         command{name: "unknown"},
			expectedErr: app.UnSubscribedCommandError{},
		},
		{
			name:    "with a subscribed command, then it execute handle method",
			cmdName: "command",
			handler: commandHandler{},
			cmd:     command{name: "command"},
		},
		{
			name:    "with a subscribed command, then it execute handle method and return same command error",
			cmdName: "other command",
			handler: commandHandler{
				err: errTest,
			},
			cmd:         command{name: "other command"},
			expectedErr: errTest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a CommandBus,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			bus := app.NewCommandBus()
			bus.Subscribe(tt.cmdName, tt.handler)
			events, err := bus.Handle(context.Background(), tt.cmd)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				require.Nil(t, events)
				return
			}
			require.Equal(t, tt.expectedEvents, events)
		})
	}
}

type command struct {
	name string
}

func (c command) Name() string {
	return c.name
}

type commandHandler struct {
	events []cqs.Event
	err    error
}

func (c commandHandler) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	return c.events, c.err
}
