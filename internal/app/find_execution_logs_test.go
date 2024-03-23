package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"

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
		expectedErr, findErr error
		logs                 []program.ExecutionLog
		expectedResult       any
		query                cqs.Query
	}{

		{
			name:        "with an invalid query, then it returns an invalid command error",
			query:       invalidQuery{},
			expectedErr: cqs.InvalidQueryError{},
		},
		{
			name:        "with an invalid limit, then it returns an invalid execution log limit error",
			query:       app.FindExecutionLogsQuery{Limit: 100},
			expectedErr: app.ErrInvalidExecutionsLogLimit,
		},
		{
			name:        "and findAll returns an error, then it returns same error",
			query:       app.FindExecutionLogsQuery{Limit: 10},
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:           "and the logs are less than limit, then it returns all logs",
			query:          app.FindExecutionLogsQuery{Limit: 10},
			logs:           lessLogs,
			expectedResult: lessLogs,
		},
		{
			name:  "and the limit is less than logs, then it returns filter logs",
			query: app.FindExecutionLogsQuery{Limit: 2},
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
			result, err := handler.Handle(context.Background(), tt.query)
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

type invalidQuery struct{}

func (i invalidQuery) Name() string {
	return "invalid"
}
