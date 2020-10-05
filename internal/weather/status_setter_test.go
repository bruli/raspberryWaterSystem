package weather_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	"github.com/stretchr/testify/assert"
)

func TestStatusSetter_Set(t *testing.T) {
	tests := map[string]struct {
		temp, hum                                float32
		rain                                     rain.Rain
		weatherErr, rainErr, expectedErr, logErr error
	}{
		"it should return error when repository return error": {
			weatherErr:  errors.New("error"),
			expectedErr: fmt.Errorf("failed reading weather data: %w", errors.New("error")),
		},
		"it should write error log when rain returns error": {
			temp:    25,
			hum:     45,
			rainErr: errors.New("error"),
			logErr:  fmt.Errorf("failed to read rain: %w", fmt.Errorf("failed to read rain data: %w", errors.New("error"))),
		},
		"it should set weather data": {
			temp: 20,
			hum:  40,
			rain: rain.New(false, 1023),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			st := status.New()
			weatherRepo := weather.RepositoryMock{}
			rainRepo := rain.RepositoryMock{}
			rainRead := rain.NewReader(&rainRepo)
			log := logger.LoggerMock{}
			s := weather.NewStatusSetter(st, &weatherRepo, rainRead, &log)

			weatherRepo.ReadFunc = func() (float32, float32, error) {
				return tt.temp, tt.hum, tt.weatherErr
			}
			rainRepo.GetFunc = func() (rain.Rain, error) {
				return tt.rain, tt.rainErr
			}
			log.FatalFunc = func(v ...interface{}) {
				assert.NotNil(t, v)
			}

			err := s.Set()
			assert.Equal(t, tt.expectedErr, err)
			if err == nil {
				assert.Equal(t, tt.temp, st.Temperature())
				assert.Equal(t, tt.hum, st.Humidity())
			}

		})
	}
}
