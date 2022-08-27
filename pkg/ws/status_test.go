package ws_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/pkg/ws"
)

func TestGetStatus(t *testing.T) {
	statusResp := http2.StatusResponseJson{
		Humidity:        10,
		IsRaining:       false,
		SystemStartedAt: vo.TimeNow().Add(-time.Hour).EpochString(),
		Temperature:     20,
		UpdatedAt:       vo.StringPtr(vo.TimeNow().EpochString()),
	}
	invalidResp := struct{ id string }{}
	started, _ := vo.ParseFromEpochStr(statusResp.SystemStartedAt)
	updated, _ := vo.ParseFromEpochStr(vo.StringValue(statusResp.UpdatedAt))
	status := ws.Status{
		Humidity:        statusResp.Humidity,
		IsRaining:       statusResp.IsRaining,
		SystemStartedAt: started,
		Temperature:     statusResp.Temperature,
		UpdatedAt:       vo.TimePtr(updated),
	}
	tests := []struct {
		name                string
		expectedErr, cliErr error
		expectedStatus      ws.Status
		response            *http.Response
	}{
		{
			name:        "and http client returns an error, then it returns a server error",
			cliErr:      errors.New(""),
			expectedErr: ws.ErrServer,
		},
		{
			name:        "and http client returns a not found response, then it returns a remote server error",
			response:    &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:        "and http client returns internal server error response, then it returns a remote server error",
			response:    &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:        "and http client returns ok but invalid response object, then it returns a failed to read response error",
			response:    &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, invalidResp)},
			expectedErr: ws.ErrFailedToReadResponse,
		},
		{
			name:           "and http client returns ok response, then it returns a valid status object",
			response:       &http.Response{StatusCode: http.StatusOK, Body: buildBody(t, statusResp)},
			expectedStatus: status,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a GetStatus function,
		when is called `+tt.name, func(t *testing.T) {
			cl := &HTTPClientMock{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tt.response, tt.cliErr
				},
			}
			pkg := ws.New(url.URL{}, cl, "token")
			st, err := pkg.GetStatus(context.Background())
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedStatus, st)
		})
	}
}

func buildBody(t *testing.T, resp interface{}) io.ReadCloser {
	data, err := json.Marshal(resp)
	require.NoError(t, err)
	body := io.NopCloser(strings.NewReader(string(data)))

	return body
}
