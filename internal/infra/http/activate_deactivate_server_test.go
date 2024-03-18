package http_test

import (
	"context"
	"errors"
	"fmt"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestActivateDeactivateServer(t *testing.T) {
	tests := []struct {
		name         string
		chErr        error
		expectedCode int
		action       string
	}{
		{
			name:         "with an invalid action name, then it returns a bad request",
			action:       "invalid",
			expectedCode: http2.StatusBadRequest,
		},
		{
			name:         "with a valid action name and command handler returns error, then it returns an internal server error",
			action:       http.ActivateAction,
			chErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "with a valid action name and command handler returns nil, then it returns ok",
			action:       http.DeactivateAction,
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given an ActivateDeactivateServer http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.ActivateDeactivateServer(ch)
			server := chi.NewMux()
			server.Patch("/status/{action}", handler)
			req := httptest.NewRequest(http2.MethodPatch, fmt.Sprintf("/status/%s", tt.action), nil)
			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)
			require.Equal(t, tt.expectedCode, writer.Code)
		})
	}
}
