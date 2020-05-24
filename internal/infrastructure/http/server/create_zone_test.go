package server

import (
	"bytes"
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateZoneHandler_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	repo := zone.RepositoryMock{}
	relays := []string{"1", "2"}
	relayRepo := zone.RelayRepositoryMock{}
	log := logger.LoggerMock{}
	router.createZone = newCreateZone(zone.NewCreator(&repo, &relayRepo, &log), &log)
	server := router.buildServer(config.AuthToken)

	tests := map[string]struct {
		body   ZoneBody
		code   int
		zones  *zone.Zones
		zone   *zone.Zone
		err    error
		relays []string
	}{
		"it should return bad request with empty body": {code: http.StatusBadRequest},
		"it should return bad request with invalid body": {
			code: http.StatusBadRequest,
			body: ZoneBody{ID: "", Name: "name"}},
		"it should return bad request with invalid relay": {
			code:   http.StatusBadRequest,
			body:   ZoneBody{ID: "aa", Name: "name", Relays: []string{"25"}},
			relays: relays},
		"it should return internal server error when save returns error": {
			code:   http.StatusInternalServerError,
			body:   ZoneBody{ID: "aa", Name: "name", Relays: []string{"1"}},
			relays: relays,
			err:    errors.New("error"),
			zones:  &zone.Zones{}},
		"it should return accepted": {
			code:   http.StatusAccepted,
			body:   ZoneBody{ID: "aa", Name: "name", Relays: []string{"1"}},
			relays: relays,
			zones:  &zone.Zones{}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			data, _ := jsoniter.Marshal(tt.body)
			request, err := http.NewRequest(http.MethodPut, "/zones", bytes.NewBuffer(data))
			request.Header.Add("Authorization", config.AuthToken)
			assert.NoError(t, err)

			repo.GetZonesFunc = func() *zone.Zones {
				return tt.zones
			}
			repo.FindFunc = func(id string) *zone.Zone {
				return tt.zone
			}
			repo.SaveFunc = func(z zone.Zones) error {
				return tt.err
			}
			relayRepo.GetFunc = func() []string {
				return tt.relays
			}
			log.FatalfFunc = func(format string, v ...interface{}) {
			}
			log.FatalFunc = func(v ...interface{}) {
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, request)
			assert.Equal(t, tt.code, writer.Code)
		})
	}
}
