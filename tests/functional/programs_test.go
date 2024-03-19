//go:build functional

package functional

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
			resp, err := buildRequestAndSend(ctx, nil, nil, http.MethodPut, "/programs", cl)
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
								Zones:   []string{savedZone.Id()},
							},
						},
						Hour: "16:13",
					},
				},
				Odd: []http2.ProgramItemRequest{
					{
						Executions: []http2.ExecutionItemRequest{
							{
								Seconds: 20,
								Zones:   []string{savedZone.Id()},
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
								Zones:   []string{savedZone.Id()},
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
										Zones:   []string{savedZone.Id()},
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
										Zones:   []string{savedZone.Id()},
									},
								},
								Hour: "08:00",
							},
						},
						Temperature: float64(25.3),
					},
				},
			}
			resp, err := buildRequestAndSend(ctx, req, authorizationHeader(), http.MethodPut, "/programs", cl)
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
