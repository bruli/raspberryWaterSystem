package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestZone(t *testing.T) {
	t.Run("it should create zone", func(t *testing.T) {
		zo, err := zone.New("bb", "bonsai big", []string{"1", "2"})
		assert.NoError(t, err)

		body := server.ZoneBody{ID: zo.Id(), Name: zo.Name(), Relays: zo.Relays()}
		data, err := jsoniter.Marshal(body)
		assert.NoError(t, err)

		resp, err := sendRequest(http.MethodPut, "/zones", data)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		respDelete, err := sendRequest(http.MethodDelete, "/zones/bb", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, respDelete.StatusCode)
	})
}
