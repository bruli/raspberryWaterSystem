package http_test

import (
	"context"
	"errors"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

func TestCreatePrograms(t *testing.T) {
	body := buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
		Daily: []http.ProgramItemRequest{
			{
				Executions: []http.ExecutionItemRequest{
					{
						Seconds: 20,
						Zones:   []string{"1"},
					},
				},
				Hour: "20:00",
			},
		},
		Odd: []http.ProgramItemRequest{
			{
				Executions: []http.ExecutionItemRequest{
					{
						Seconds: 20,
						Zones:   []string{"1"},
					},
				},
				Hour: "20:00",
			},
		},
		Even: []http.ProgramItemRequest{
			{
				Executions: []http.ExecutionItemRequest{
					{
						Seconds: 20,
						Zones:   []string{"1"},
					},
				},
				Hour: "20:00",
			},
		},
		Weekly: []http.WeeklyItemRequest{
			{
				Programs: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"4"},
							},
						},
						Hour: "15:10",
					},
				},
				WeekDay: "Friday",
			},
		},
		Temperature: []http.TemperatureItemRequest{
			{
				Programs: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"1"},
							},
						},
						Hour: "08:00",
					},
				},
				Temperature: 35,
			},
		},
	})
	tests := []struct {
		name, body   string
		chErr        error
		expectedCode int
	}{
		{
			name:         "with an invalid request, then it returns bad request",
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid program hour, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Daily: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"1"},
							},
						},
						Hour: "invalid",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid program seconds, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Daily: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{
							{
								Seconds: -20,
								Zones:   []string{"1"},
							},
						},
						Hour: "20:00",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid program zones, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Daily: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{},
							},
						},
						Hour: "20:00",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid program executions, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Daily: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{},
						Hour:       "20:00",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid odd program, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Odd: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{},
						Hour:       "20:00",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid even program, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Even: []http.ProgramItemRequest{
					{
						Executions: []http.ExecutionItemRequest{},
						Hour:       "20:00",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid weekday, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Weekly: []http.WeeklyItemRequest{
					{
						Programs: []http.ProgramItemRequest{},
						WeekDay:  "invalid",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid weekly program, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Weekly: []http.WeeklyItemRequest{
					{
						Programs: []http.ProgramItemRequest{
							{
								Executions: []http.ExecutionItemRequest{},
								Hour:       "invalid",
							},
						},
						WeekDay: "Friday",
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid temperature programs, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Temperature: []http.TemperatureItemRequest{
					{
						Programs: []http.ProgramItemRequest{
							{
								Executions: []http.ExecutionItemRequest{},
								Hour:       "invalid",
							},
						},
						Temperature: 0,
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name: "with a valid request with invalid temperature value, then it returns bad request",
			body: buildRequestJsonToString(t, &http.CreateProgramsRequestJson{
				Temperature: []http.TemperatureItemRequest{
					{
						Programs: []http.ProgramItemRequest{
							{
								Executions: []http.ExecutionItemRequest{
									{
										Seconds: 20,
										Zones:   []string{"1"},
									},
								},
								Hour: "15:00",
							},
						},
						Temperature: float64(-15),
					},
				},
			}),
			expectedCode: http2.StatusBadRequest,
		},
		{
			name:         "with a valid request and command handler returns error, then it returns bad request",
			body:         body,
			chErr:        errors.New(""),
			expectedCode: http2.StatusInternalServerError,
		},
		{
			name:         "with a valid request, then it returns ok",
			body:         body,
			expectedCode: http2.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a create programs http handler,
		when a request is sent `+tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
					return nil, tt.chErr
				},
			}
			handler := http.CreatePrograms(ch)
			req := httptest.NewRequest(http2.MethodPost, "/programs", buildRequestBody(tt.body))
			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, req)
			resp := writer.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
