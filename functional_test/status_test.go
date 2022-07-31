//go:build functional
// +build functional

package functional_test

import (
	"net/http"
	"testing"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"

	"github.com/stretchr/testify/require"
)

func runStatus(t *testing.T) {
	t.Run(`Given a FindStatus endpoint,
	when a request is sent `, func(t *testing.T) {
		t.Run(`without authorization, then it return unauthorized`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, nil, http.MethodGet, "/status", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`wit authorization, then it return a valid response`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, authorizationHeader(), http.MethodGet, "/status", cl)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var schema http2.StatusResponseJson
			readResponse(t, resp, &schema)
			require.NotEmpty(t, schema.SystemStartedAt)
		})
	})
}
