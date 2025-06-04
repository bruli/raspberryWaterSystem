package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRemoveZone(t *testing.T) {
	tests := []struct {
		name, body   string
		chErr        error
		expectedCode int
	}{
		{
			name:         "and command handler return a not found error, then it returns a not found",
			chErr:        vo.NotFoundError{},
			expectedCode: http2.StatusNotFound,
		},
		{
			name:         "and command handler return an error, then it returns an internal server error",
			chErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "and command handler return nil, then it returns ok",
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Remove zone http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.RemoveZone(ch)
			server := chi.NewMux()
			server.Delete("/zones/{id}", handler)
			req := httptest.NewRequest(http2.MethodDelete, "/zones/id", buildRequestBody(tt.body))
			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
