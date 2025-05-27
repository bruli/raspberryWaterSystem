package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

func TestExecuteZone(t *testing.T) {
	body := buildRequestJsonToString(t, http.ExecuteZoneRequestJson{Seconds: 20})
	tests := []struct {
		name, body   string
		chErr        error
		expectedCode int
	}{
		{
			name:         "with invalid request, then it returns bad request",
			expectedCode: http2.StatusBadRequest,
			body:         "invalid",
		},
		{
			name:         "with a valid request and command handler returns not found error, then it returns not found",
			expectedCode: http2.StatusNotFound,
			body:         body,
			chErr:        vo.NotFoundError{},
		},
		{
			name:         "with a valid request and command handler returns en error, then it returns internal server error",
			expectedCode: http2.StatusInternalServerError,
			body:         body,
			chErr:        errors.New(""),
		},
		{
			name:         "with a valid request, then it returns ok",
			expectedCode: http2.StatusOK,
			body:         body,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a ExecuteZone http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.ExecuteZone(ch)
			server := chi.NewMux()
			server.Post("/zones/{id}/execute", handler)
			req := httptest.NewRequest(http2.MethodPost, "/zones/id/execute", buildRequestBody(tt.body))
			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
