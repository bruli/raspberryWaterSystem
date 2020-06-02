package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getRain_ServeHTTP(t *testing.T) {
	config := getConfig()
	router := getRouter()
	log := logger.LoggerMock{}
	repo := rain.RepositoryMock{}
	reader := rain.NewReader(&repo)
	st := status.New()
	router.getRain = newGetRain(reader, &log, st)
	server := router.buildServer(config.AuthToken)
	data := rain.New(true, 200)
	tests := map[string]struct {
		rainData rain.Rain
		code     int
		err      error
	}{
		"it should return internal server error when reader returns error": {
			rainData: data,
			code:     http.StatusInternalServerError,
			err:      errors.New("error"),
		},
		"it should return rain data": {
			rainData: data,
			code:     http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/rain", nil)
			request.Header.Add("Authorization", config.AuthToken)
			assert.NoError(t, err)

			repo.GetFunc = func() (rain.Rain, error) {
				return tt.rainData, tt.err
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, request)
			assert.Equal(t, tt.code, writer.Code)
			if tt.err == nil {
				assert.Equal(t, st.Rain().Value(), data.Value())
				assert.Equal(t, st.Rain().IsRain(), data.IsRain())

				body := rainBody{}
				err = jsoniter.Unmarshal(writer.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, tt.rainData.Value(), body.Value)
				assert.Equal(t, tt.rainData.IsRain(), body.IsRain)
			}
		})
	}
}
