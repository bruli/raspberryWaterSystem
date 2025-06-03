package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestFindWeather(t *testing.T) {
	weath := fixtures.WeatherBuilder{}.Build()
	tests := []struct {
		name         string
		expectedCode int
		qhErr        error
		result       any
	}{
		{
			name:         "and query handler returns an error, then it returns an internal server error",
			qhErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "and query handler returns a result, then it returns a valid result",
			result:       weath,
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {

		t.Run(`Given a FindWeather http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			qh := &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqs.Query) (any, error) {
					return tt.result, tt.qhErr
				},
			}
			handler := http.FindWeather(qh)
			req := httptest.NewRequest(http2.MethodGet, "/weather", nil)
			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
			if resp.StatusCode == http2.StatusOK {
				var schema http.WeatherResponseJson
				readResponse(t, resp, &schema)
			}
		})
	}
}
