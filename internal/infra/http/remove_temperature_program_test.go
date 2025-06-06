package http_test

import (
	"context"
	"errors"
	"fmt"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRemoveTemperatureProgram(t *testing.T) {
	temp := "23"
	tests := []struct {
		name         string
		temp         string
		chErr        error
		expectedCode int
	}{
		{
			name:         "with an invalid temp, then it returns a bad request",
			temp:         "invalid",
			expectedCode: http2.StatusBadRequest,
		},
		{
			name:         "and command handler returns a not found error, then it returns a not found",
			temp:         temp,
			chErr:        vo.NotFoundError{},
			expectedCode: http2.StatusNotFound,
		},
		{
			name:         "and command handler returns an error, then it returns an internal server error",
			temp:         temp,
			chErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "and command handler returns nil, then it returns ok status",
			temp:         temp,
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a RemoveTemperatureProgram http handler,
		when a request is sent`+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.RemoveTemperatureProgram(ch)
			server := chi.NewMux()
			server.Delete("/programs/temperature/{temperature}", handler)
			req := httptest.NewRequest(http2.MethodDelete, fmt.Sprintf("/programs/temperature/%s", tt.temp), nil)
			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
