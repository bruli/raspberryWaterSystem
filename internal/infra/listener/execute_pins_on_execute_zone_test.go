package listener_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/internal/infra/listener"
)

func TestExecutePinsOnExecuteZoneListen(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name string
		expectedErr, executeCHErr,
		logsCHErr, publishCHErr error
	}{
		{
			name:         "and execute pins command handler returns an error, then it returns same error",
			executeCHErr: errTest,
			expectedErr:  errTest,
		},
		{
			name:        "and save execution logs command handler returns an error, then it returns same error",
			logsCHErr:   errTest,
			expectedErr: errTest,
		},
		{
			name:         "and publish message command handler returns an error, then it returns same error",
			publishCHErr: errTest,
			expectedErr:  errTest,
		},
		{
			name: "and all services works fine, then it returns nil",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a ExecutePinsOnExecuteZone listener,
		when Listen method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{}
			ch.HandleFunc = func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
				_, isExecute := cmd.(app.ExecutePinsCmd)
				_, isLog := cmd.(app.SaveExecutionLogCmd)
				_, isPublish := cmd.(app.PublishMessageCmd)
				switch {
				case isExecute:
					return nil, tt.executeCHErr
				case isLog:
					return nil, tt.logsCHErr
				case isPublish:
					return nil, tt.publishCHErr
				default:
					return nil, nil
				}
			}
			list := listener.NewExecutePinsOnExecuteZone(ch)
			err := list.Listen(context.Background(), zone.Executed{
				ZoneName: "zone test",
				Seconds:  10,
			})
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
