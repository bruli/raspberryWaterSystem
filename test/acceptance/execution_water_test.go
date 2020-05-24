package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestExecutionWater(t *testing.T) {
	body := server.NewExecuteWaterBody(1, "bb")
	data, err := jsoniter.Marshal(body)
	assert.NoError(t, err)
	resp, err := sendRequest(http.MethodPost, "/water", data)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}
