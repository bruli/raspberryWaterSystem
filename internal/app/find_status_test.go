package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/stretchr/testify/require"
)

func TestFindStatusHandle(t *testing.T) {
	errTest := errors.New("")
	st := fixtures.StatusBuilder{}.Build()
	tests := []struct {
		name                 string
		expectedErr, findErr error
		status               status.Status
		expectedResult       any
	}{
		{
			name:        "and find method returns an error, then it returns same error",
			findErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:           "and find method returns a status, then it returns a valid result",
			status:         st,
			expectedResult: st,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a FindStatus query handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			sr := &StatusRepositoryMock{
				FindFunc: func(ctx context.Context) (status.Status, error) {
					return tt.status, tt.findErr
				},
			}
			handler := app.NewFindStatus(sr)
			result, err := handler.Handle(context.Background(), app.FindStatusQuery{})
			if err != nil {
				test.CheckErrorsType(t, tt.expectedErr, err)
				return
			}
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
