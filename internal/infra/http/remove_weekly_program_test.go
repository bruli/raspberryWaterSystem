package http_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	http2 "net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoveWeeklyProgram(t *testing.T) {
	day := "monday"
	tests := []struct {
		name         string
		day          string
		chErr        error
		expectedCode int
	}{
		{
			name:         "with an invalid day, then it returns a bad request",
			day:          "invalid",
			expectedCode: http2.StatusBadRequest,
		},
		{
			name:         "and command handler returns a not found error, then it returns a not found",
			day:          day,
			chErr:        vo.NotFoundError{},
			expectedCode: http2.StatusNotFound,
		},
		{
			name:         "and command handler returns an error, then it returns an internal server error",
			day:          day,
			chErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "and command handler returns nil, then it returns ok status",
			day:          day,
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a RemoveWeeklyProgram http handler,
		when a request is sent`+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.RemoveWeeklyProgram(ch)
			server := chi.NewMux()
			server.Delete("/programs/weekly/{day}", handler)
			req := httptest.NewRequest(http2.MethodDelete, fmt.Sprintf("/programs/weekly/%s", tt.day), nil)
			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
