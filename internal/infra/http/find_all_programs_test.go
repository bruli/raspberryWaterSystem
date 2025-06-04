package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/stretchr/testify/require"
)

func TestFindAllPrograms(t *testing.T) {
	tests := []struct {
		name         string
		result       any
		qhErr        error
		expectedCode int
	}{
		{
			name:         "and query handler return an error, then it returns an internal server error",
			qhErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name: "and query handler return a result, then it returns a valid response",
			result: app.AllPrograms{
				Daily: []program.Program{
					fixtures.ProgramBuilder{}.Build(),
				},
				Odd: []program.Program{
					fixtures.ProgramBuilder{}.Build(),
				},
				Even: []program.Program{
					fixtures.ProgramBuilder{}.Build(),
				},
				Weekly: []program.Weekly{
					fixtures.WeeklyBuilder{}.Build(),
				},
				Temperature: []program.Temperature{
					fixtures.TemperatureBuilder{}.Build(),
				},
			},
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a FindAllPrograms http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			qh := &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqs.Query) (any, error) {
					return tt.result, tt.qhErr
				},
			}
			handler := http.FindAllPrograms(qh)
			req := httptest.NewRequest(http2.MethodGet, "/programs", nil)
			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
			if resp.StatusCode == http2.StatusOK {
				var schema http.ProgramsResponseJson
				readResponse(t, resp, &schema)
			}
		})
	}
}
