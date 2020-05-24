package acceptance

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTemperature(t *testing.T) {
	resp, err := sendRequest(http.MethodGet, "/temperature", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
