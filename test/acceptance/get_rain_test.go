package acceptance

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRain(t *testing.T) {
	resp, err := sendRequest(http.MethodGet, "/rain", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
