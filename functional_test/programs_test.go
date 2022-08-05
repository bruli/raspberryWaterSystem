//go:build functional
// +build functional

package functional_test

import (
	"net/http"
	"testing"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/stretchr/testify/require"
)

func runPrograms(t *testing.T) {
	t.Run(`Given a create programs endpoint,
	when a request is sent`, func(t *testing.T) {
		t.Run(`without authorization,
		then it returns unauthorized`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, nil, http.MethodPost, "/programs", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`with authorization,
		then it returns unauthorized`, func(t *testing.T) {
			req := http2.CreateProgramsRequestJson{
				Daily: []http2.ProgramItemRequest{
					{
						Executions: []http2.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"d8aa59a8-6ce2-4cdd-95e5-ca03adfaec67", "dafc349a-a5bf-413f-a59a-0a6a71f095d9"},
							},
						},
						Hour: "08:07",
					},
				},
				Odd: []http2.ProgramItemRequest{
					{
						Executions: []http2.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"1"},
							},
						},
						Hour: "15:10",
					},
				},
				Even: []http2.ProgramItemRequest{
					{
						Executions: []http2.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{"1"},
							},
						},
						Hour: "15:10",
					},
				},
				Weekly: []http2.WeeklyItemRequest{
					{
						Programs: []http2.ProgramItemRequest{
							{
								Executions: []http2.ExecutionItemRequest{
									{
										Seconds: 15,
										Zones:   []string{"1", "2"},
									},
								},
								Hour: "08:00",
							},
						},
						WeekDay: "Friday",
					},
				},
				Temperature: []http2.TemperatureItemRequest{
					{
						Programs: []http2.ProgramItemRequest{
							{
								Executions: []http2.ExecutionItemRequest{
									{
										Seconds: 15,
										Zones:   []string{"1", "2"},
									},
								},
								Hour: "08:00",
							},
						},
						Temperature: float64(25.3),
					},
				},
			}
			resp, err := buildRequestAndSend(ctx, req, authorizationHeader(), http.MethodPost, "/programs", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
	t.Run(`Given a find all programs endpoint,
	when a request is sent`, func(t *testing.T) {
		t.Run(`without authorization,
		then it returns an unauthorized`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, nil, http.MethodGet, "/programs", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`with authorization,
		then it returns a valid response`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, authorizationHeader(), http.MethodGet, "/programs", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var schema http2.ProgramsResponseJson
			readResponse(t, resp, &schema)
		})
	})
}
