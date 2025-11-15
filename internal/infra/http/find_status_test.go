package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/fixtures"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"

	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

func TestFindStatus(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		result       any
		qhErr        error
	}{
		{
			name:         "and query handler returns a not found error, then it returns a not found",
			qhErr:        vo.NotFoundError{},
			expectedCode: http2.StatusNotFound,
		},
		{
			name:         "and query handler returns an error, then it returns an internal server error",
			qhErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "and query handler returns an result, then it returns a valid response",
			result:       fixtures.StatusBuilder{}.Build(),
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a FindStatus query handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			qh := &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqs.Query) (any, error) {
					return tt.result, tt.qhErr
				},
			}
			handler := http.FindStatus(qh)
			req := httptest.NewRequest(http2.MethodGet, "/status", nil)
			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
			if resp.StatusCode == http2.StatusOK {
				var schema http.StatusResponseJson
				readResponse(t, resp, &schema)
			}
		})
	}
}
