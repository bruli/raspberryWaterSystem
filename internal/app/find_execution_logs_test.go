package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestFindExecutionLogsHandle(t *testing.T) {
	errTest := errors.New("")
	executedTime1 := vo.TimeNow()
	executedTime2 := executedTime1.AddDate(0, 0, -1)
	executedTime3 := executedTime1.AddDate(0, 0, -2)
	firstLog := fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 1"), ExecutedAt: &executedTime1}.Build()
	secondLog := fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 2"), ExecutedAt: &executedTime2}.Build()
	thirdLog := fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 3"), ExecutedAt: &executedTime3}.Build()
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
			name:  "and the logs are less than limit, then it returns all logs",
			query: app.FindExecutionLogsQuery{Limit: 10},
			logs: []program.ExecutionLog{
				thirdLog,
				secondLog,
				firstLog,
			},
			expectedResult: []program.ExecutionLog{
				firstLog,
				secondLog,
				thirdLog,
			},
		},
		{
			name:  "and the limit is less than logs, then it returns filter logs",
			query: app.FindExecutionLogsQuery{Limit: 2},
			logs: []program.ExecutionLog{
				thirdLog,
				secondLog,
				firstLog,
			},
			expectedResult: []program.ExecutionLog{
				firstLog,
				secondLog,
			},
		},
	}
	for _, tt := range tests {
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
				require.ErrorAs(t, err, &tt.expectedErr)
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
