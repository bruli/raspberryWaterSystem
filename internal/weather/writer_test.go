package weather

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	tests := map[string]struct {
		temp, hum                      float32
		expectedErr, readErr, writeErr error
	}{
		"it should return error when repository returns error": {
			readErr:     errors.New("error"),
			expectedErr: fmt.Errorf("failed to read weather data: %w", errors.New("error")),
		},
		"it should return error when writer returns error": {
			temp:        22,
			hum:         40,
			writeErr:    errors.New("error"),
			expectedErr: fmt.Errorf("failed to write weather data: %w", errors.New("error")),
		},
		"it should write weather": {
			temp: 22,
			hum:  40,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := RepositoryMock{}
			writeRepo := WriteRepositoryMock{}
			writ := NewWriter(&repo, &writeRepo)

			repo.ReadFunc = func() (float32, float32, error) {
				return tt.temp, tt.hum, tt.readErr
			}
			writeRepo.WriteFunc = func(temp float32, hum float32) error {
				return tt.writeErr
			}

			err := writ.Write()
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
