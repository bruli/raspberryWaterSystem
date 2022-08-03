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
