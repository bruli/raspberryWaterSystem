package http_test

import (
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name, token  string
		expectedCode int
	}{
		{
			name:         "without Authorization token, then it returns an unauthorized",
			expectedCode: http2.StatusUnauthorized,
		},
		{
			name:         "with invalid Authorization token, then it returns an unauthorized",
			token:        "invalid",
			expectedCode: http2.StatusUnauthorized,
		},
		{
			name:         "with valid Authorization token, then it returns an ok",
			token:        "1234",
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given an AuthMiddleware,
		when is used and a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			middleware := http.AuthMiddleware("1234")
			handler := middleware(nextHandler())
			req := httptest.NewRequest(http2.MethodPost, "/", nil)
			if len(tt.token) != 0 {
				req.Header.Set("Authorization", tt.token)
			}
			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func nextHandler() http2.HandlerFunc {
	return func(w http2.ResponseWriter, r *http2.Request) {
		httpx.WriteResponse(w, http2.StatusOK, nil)
	}
}
