package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetZonesHandler_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	logger := logger.LoggerMock{}
	repo := zone.RepositoryMock{}
	router.getZones = newGetZones(zone.NewGetter(&repo), &logger)
	server := router.buildServer(config.AuthToken)
	zo, err := zone.NewZonesStub()
	assert.NoError(t, err)

	tests := map[string]struct {
		zones zone.Zones
	}{
		"it should return empty zones": {zones: zone.Zones{}},
		"it should return zones":       {zones: zo},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/zones", nil)
			request.Header.Add("Authorization", config.AuthToken)
			assert.NoError(t, err)

			repo.GetZonesFunc = func() *zone.Zones {
				return &tt.zones
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, request)
			assert.Equal(t, http.StatusOK, writer.Code)

			body := ZonesBody{}
			err = jsoniter.Unmarshal(writer.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.zones), len(body))
		})
	}
}
