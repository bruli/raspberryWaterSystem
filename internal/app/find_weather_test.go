package app_test

import (
	"context"
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/test"
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
		expectedResult cqs.QueryResult
		temp, hum      float32
		rain           bool
	}{
		{
			name:        "and find temperature returns an error, then it returns same error",
			tempErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:        "and find rain returns an error, then it returns same error",
			rainErr:     errTest,
			expectedErr: errTest,
			temp:        temp,
			hum:         hum,
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
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
