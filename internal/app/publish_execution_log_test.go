package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/stretchr/testify/require"
)

func TestPublishExecutionLogHandle(t *testing.T) {
	errTest := errors.New("")
	cmd := app.PublishExecutionLogCmd{
		ZoneName:   "zone",
		Seconds:    program.Seconds(time.Second),
		ExecutedAt: vo.TimeNow(),
	}
	tests := []struct {
		name                    string
		publishErr, expectedErr error
		cmd                     cqs.Command
	}{
		{
			name:        "and create execution log returns an error, then it returns a publish execution log error",
			cmd:         app.PublishExecutionLogCmd{},
			expectedErr: app.PublishExecutionLogError{},
		},
		{
			name:        "and publish returns an error, then it returns same error",
			cmd:         cmd,
			publishErr:  errTest,
			expectedErr: errTest,
		},
		{
			name: "then it returns nil",
			cmd:  cmd,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a PublishExecutionLog command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			elp := &ExecutionLogPublisherMock{
				PublishFunc: func(ctx context.Context, execLog program.ExecutionLog) error {
					return tt.publishErr
				},
			}
			handler := app.NewPublishExecutionLog(elp)
			events, err := handler.Handle(context.Background(), tt.cmd)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Nil(t, events)
		})
	}
}
