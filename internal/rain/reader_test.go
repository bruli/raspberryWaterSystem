package rain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	tests := map[string]struct {
		r                rain.Rain
		err, expectedErr error
	}{
		"it should return error when repository returns error": {r: rain.Rain{},
			err:         errors.New("error"),
			expectedErr: fmt.Errorf("failed to read rain data: %w", errors.New("error"))},
		"it should return rain data": {r: rain.New(true, 200)},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := rain.RepositoryMock{}
			repo.GetFunc = func() (rain.Rain, error) {
				return tt.r, tt.err
			}
			read := rain.NewReader(&repo)
			r, err := read.Read()
			assert.Equal(t, tt.r, r)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
