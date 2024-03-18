//go:build functional

package functional_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func runHomepage(t *testing.T) {
	t.Run(`Given a homepage endpoint,
	when a request is sent,
	then it returns an OK`, func(t *testing.T) {
		resp, err := buildRequestAndSend(ctx, nil, nil, http.MethodGet, "/", cl)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
