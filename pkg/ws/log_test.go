package ws_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"

	"github.com/bruli/raspberryWaterSystem/pkg/ws"
	"github.com/stretchr/testify/require"
)

func TestGetLog(t *testing.T) {
	errorResponse := httpx.ErrorResponseJson{Errors: []httpx.Error{
		{
			Code:    "bbb",
			Message: "invalid",
		},
	}}
	logsDomain := []program.ExecutionLog{
		fixtures.ExecutionLogBuilder{}.Build(),
		fixtures.ExecutionLogBuilder{}.Build(),
	}
	logsResponse := make([]http2.ExecutionLogItemResponse, len(logsDomain))
	logs := make([]ws.Log, len(logsDomain))
	for i, lo := range logsDomain {
		logsResponse[i] = http2.ExecutionLogItemResponse{
			ExecutedAt: lo.ExecutedAt().EpochString(),
			Seconds:    lo.Seconds().Int(),
			ZoneName:   lo.ZoneName(),
		}
		logs[i] = ws.Log{
			ExecutedAt: lo.ExecutedAt(),
			Seconds:    lo.Seconds().Int(),
			ZoneName:   lo.ZoneName(),
		}
	}
	tests := []struct {
		name                string
		cliErr, expectedErr error
		response            *http.Response
		expectedLogs        []ws.Log
	}{
		{
			name:        "and http client returns an error, then it returns a server error",
			cliErr:      errors.New(""),
			expectedErr: ws.ErrServer,
		},
		{
			name:        "and http client returns internal server error response, then it returns a remote server error",
			response:    &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:        "and http client returns bad request response, then it returns a log error",
			response:    &http.Response{StatusCode: http.StatusBadRequest, Body: buildBody(t, errorResponse)},
			expectedErr: ws.LogError{},
		},
		{
			name:        "and http client returns ok with invalid response schema, then it returns failed to read response error",
			response:    &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, invalidResponse{})},
			expectedErr: ws.ErrFailedToReadResponse,
		},
		{
			name:         "and http client returns ok response, then it returns logs",
			response:     &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, logsResponse)},
			expectedLogs: logs,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a GetLog function,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			cl := &HTTPClientMock{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tt.response, tt.cliErr
				},
			}
			pkg := ws.New(url.URL{}, cl, "token")
			result, err := pkg.GetLogs(context.Background(), 2)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, len(tt.expectedLogs), len(result))
		})
	}
}
