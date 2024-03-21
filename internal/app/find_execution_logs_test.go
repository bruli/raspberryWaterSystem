package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestFindExecutionLogsHandle(t *testing.T) {
	errTest := errors.New("")
	lessLogs := []program.ExecutionLog{
		fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 1")}.Build(),
		fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 2")}.Build(),
		fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 3")}.Build(),
	}
	tests := []struct {
		name                 string
		limit                int
		expectedErr, findErr error
		logs                 []program.ExecutionLog
		expectedResult       any
	}{
		{
			name:        "with an invalid limit, then it returns an invalid execution log limit error",
			limit:       100,
			expectedErr: app.ErrInvalidExecutionsLogLimit,
		},
		{
			name:        "and findAll returns an error, then it returns same error",
			limit:       10,
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:           "and the logs are less than limit, then it returns all logs",
			limit:          10,
			logs:           lessLogs,
			expectedResult: lessLogs,
		},
		{
			name:  "and the limit is less than logs, then it returns filter logs",
			limit: 2,
			logs:  lessLogs,
			expectedResult: []program.ExecutionLog{
				lessLogs[1],
				lessLogs[2],
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a FindExecutionLogs query handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			elr := &ExecutionLogRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]program.ExecutionLog, error) {
					return tt.logs, tt.findErr
				},
			}
			handler := app.NewFindExecutionLogs(elr)
			result, err := handler.Handle(context.Background(), app.FindExecutionLogsQuery{Limit: tt.limit})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
