package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/stretchr/testify/require"
)

func TestFindWeatherHandle(t *testing.T) {
	errTest := errors.New("")
	temp := float32(23.4)
	hum := float32(45.7)
	weath := weather.New(temp, hum, false)
	tests := []struct {
		name string
		expectedErr, tempErr,
		rainErr error
		expectedResult any
		temp, hum      float32
		rain           bool
	}{
		{
			name:           "and find temperature returns an error, then it returns zero temperature",
			tempErr:        errTest,
			expectedResult: weather.New(0, 0, false),
		},
		{
			name:           "and find rain returns an error, then it returns zero values",
			rainErr:        errTest,
			temp:           temp,
			hum:            hum,
			expectedResult: weather.New(temp, hum, false),
		},
		{
			name:           "then it returns a valid result",
			temp:           temp,
			hum:            hum,
			expectedResult: weath,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a FindWeather query handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			tr := &TemperatureRepositoryMock{
				FindFunc: func(ctx context.Context) (float32, float32, error) {
					return tt.temp, tt.hum, tt.tempErr
				},
			}
			rr := &RainRepositoryMock{
				FindFunc: func(ctx context.Context) (bool, error) {
					return false, tt.rainErr
				},
			}
			handler := app.NewFindWeather(tr, rr)
			result, err := handler.Handle(context.Background(), app.FindWeatherQuery{})
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
