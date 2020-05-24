package server

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_removeZone_ServeHTTP(t *testing.T) {
	zons := zone.Zones{}
	zon, err := zone.New("1", "test", []string{"1"})
	assert.NoError(t, err)
	zons.Add(*zon)
	config := getConfig()
	router := getRouter()
	repo := zone.RepositoryMock{}
	log := logger.LoggerMock{}
	remover := zone.NewRemover(&repo, &log)
	router.removeZone = newRemoveZone(remover, &log)
	server := router.buildServer(config.AuthToken)
	tests := map[string]struct {
		zoneID string
		code   int
		zon    *zone.Zone
		zons   *zone.Zones
	}{
		"it should return not found when zone does not exists": {
			zoneID: "1",
			code:   http.StatusNotFound,
		},
		"it should return accepted": {
			zoneID: "1",
			zon:    zon,
			zons:   &zons,
			code:   http.StatusAccepted,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/zones/%s", tt.zoneID), nil)
			assert.NoError(t, err)
			request.Header.Add("Authorization", config.AuthToken)

			repo.FindFunc = func(id string) *zone.Zone {
				return tt.zon
			}
			log.FatalFunc = func(v ...interface{}) {
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}
			repo.GetZonesFunc = func() *zone.Zones {
				return tt.zons
			}
			repo.SaveFunc = func(z zone.Zones) error {
				return nil
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, request)
			assert.Equal(t, tt.code, writer.Code)
		})
	}
}
