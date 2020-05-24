package execution_test

import (
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	zon, err := zone.NewZoneStub()
	assert.NoError(t, err)
	tests := map[string]struct {
		seconds                                                                  uint8
		zoneID                                                                   string
		zone                                                                     *zone.Zone
		activateErr, deactivateErr, logRepoErr, notificationSendErr, expectedErr error
	}{
		"it should return error with invalid seconds": {
			seconds:     0,
			expectedErr: execution.NewInvalidExecutorData("invalid seconds to execute"),
		},
		"it should return error with invalid zone": {
			seconds:     1,
			expectedErr: execution.NewInvalidExecutorData("invalid zone to execute"),
		},
		"it should return error when zone repository returns nil": {
			seconds:     1,
			zoneID:      "1",
			expectedErr: errors.New("'1' is not a valid zone id"),
		},
		"it should return error when activate returns error": {
			seconds:     1,
			zoneID:      "1",
			zone:        &zon,
			activateErr: errors.New("error"),
			expectedErr: fmt.Errorf("failed activating pins: %w", errors.New("error")),
		},
		"it should return error when deactivate returns error": {
			seconds:       1,
			zoneID:        "1",
			zone:          &zon,
			deactivateErr: errors.New("error"),
			expectedErr:   fmt.Errorf("failed deactivating pins: %w", errors.New("error")),
		},
		"it should return error when log repository returns error": {
			seconds:     1,
			zoneID:      "1",
			zone:        &zon,
			logRepoErr:  errors.New("error"),
			expectedErr: fmt.Errorf("failed saving execution log: %w", errors.New("error")),
		},
		"it should return error when log notification sender returns error": {
			seconds:             1,
			zoneID:              "1",
			zone:                &zon,
			notificationSendErr: errors.New("error"),
			expectedErr:         fmt.Errorf("failed sending execution notification %w", errors.New("error")),
		},
		"it should execute": {
			seconds: 1,
			zoneID:  "1",
			zone:    &zon,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			zoneRepo := zone.RepositoryMock{}
			relayManang := execution.RelayManagerMock{}
			logRepo := execution.LogRepositoryMock{}
			notificationSend := execution.NotificationSenderMock{}
			exec := execution.NewExecutor(&zoneRepo, &relayManang, &logRepo, &notificationSend)
			zoneRepo.FindFunc = func(id string) *zone.Zone {
				return tt.zone
			}
			relayManang.ActivatePinsFunc = func(pins []string) error {
				return tt.activateErr
			}
			relayManang.DeactivatePinsFunc = func(pins []string) error {
				return tt.deactivateErr
			}
			logRepo.SaveFunc = func(l execution.Log) error {
				return tt.logRepoErr
			}
			notificationSend.SendFunc = func(message string) error {
				return tt.notificationSendErr
			}
			err := exec.Execute(tt.seconds, tt.zoneID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
