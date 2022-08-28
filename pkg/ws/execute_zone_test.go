package ws_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/ws"
	"github.com/stretchr/testify/require"
)

func TestExecuteZone(t *testing.T) {
	tests := []struct {
		name                string
		response            *http.Response
		cliErr, expectedErr error
	}{
		{
			name:        "and http client returns an error, then it returns a server error",
			cliErr:      errors.New(""),
			expectedErr: ws.ErrServer,
		},
		{
			name:        "and http client returns a not found response, then it returns an unknown zone to execute error",
			response:    &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody},
			expectedErr: ws.ErrUnknownZoneToExecute,
		},
		{
			name:        "and http client returns an internal server error response, then it returns a remote server error",
			response:    &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody},
			expectedErr: ws.ErrRemoteServerErr,
		},
		{
			name:     "and http client returns an ok response, then it returns nil",
			response: &http.Response{StatusCode: http.StatusOK, Body: http.NoBody},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a ExecuteZone method,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			cl := &HTTPClientMock{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tt.response, tt.cliErr
				},
			}
			pkg := ws.New(url.URL{}, cl, "token")
			err := pkg.ExecuteZone(context.Background(), "bbf", 2)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
