package acceptance

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestExecutionLogs(t *testing.T) {
	resp, err := sendRequest(http.MethodGet, "/executions/logs", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
