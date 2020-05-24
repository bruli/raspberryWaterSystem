package acceptance

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetExecutions(t *testing.T) {
	t.Run("it should return executions", func(t *testing.T) {
		resp, err := sendRequest(http.MethodGet, "/executions", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
