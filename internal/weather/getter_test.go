package weather_test

import (
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetter_Get(t *testing.T) {
	tests := map[string]struct {
		temp, hum        float32
		err, expectedErr error
	}{
		"it should return error when repository returns error": {
			temp:        0,
			hum:         0,
			err:         errors.New("error"),
			expectedErr: fmt.Errorf("failed reading weather data: %w", errors.New("error")),
		},
		"it should return weather data": {
			temp: 20,
			hum:  40,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := weather.RepositoryMock{}
			get := weather.NewGetter(&repo)

			repo.ReadFunc = func() (float32, float32, error) {
				return tt.temp, tt.hum, tt.err
			}

			temp, hum, err := get.Get()

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.temp, temp)
			assert.Equal(t, tt.hum, hum)
		})
	}
}
