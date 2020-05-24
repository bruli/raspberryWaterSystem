package execution_test

import (
	"errors"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecutorInTime_Execute(t *testing.T) {
	now := time.Now()
	hour := now.Format("15:04")
	daily := execution.Programs{}
	prgm, err := execution.NewProgram(1, hour, []string{"1", "2"})
	assert.NoError(t, err)
	daily.Add(prgm)
	week := execution.WeeklyPrograms{}
	odd := execution.Programs{}
	even := execution.Programs{}
	prgms, err := execution.New(&daily, &week, &odd, &even)
	assert.NoError(t, err)
	zon, err := zone.New("1", "name1", []string{"a"})
	assert.NoError(t, err)
	tests := map[string]struct {
		t                                   time.Time
		expectedErr, repositoryErr, sendErr error
		rain                                rain.Rain
		execut                              *execution.Execution
		zon                                 *zone.Zone
	}{
		"it should return error when repository returns error": {
			t:             now,
			repositoryErr: errors.New("error"),
			expectedErr:   fmt.Errorf("failed to get executions: %w", errors.New("error")),
		},
		"it should return error when is raining and notification sender returns error": {
			t:           now,
			execut:      prgms,
			rain:        rain.New(true, 200),
			sendErr:     errors.New("error"),
			expectedErr: fmt.Errorf("failed to send execution notification: %w", errors.New("error")),
		},
		"it should return error when executor returns error": {
			t:           now,
			execut:      prgms,
			rain:        rain.New(false, 1023),
			expectedErr: fmt.Errorf("failed to execute zone '%s': %w", "1", fmt.Errorf("'%s' is not a valid zone id", "1")),
		},
		"it should execute": {
			t:      now,
			execut: prgms,
			zon:    zon,
			rain:   rain.New(false, 1023),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := execution.RepositoryMock{}
			st := status.New()
			st.SetRain(tt.rain.IsRain(), tt.rain.Value())
			sender := execution.NotificationSenderMock{}
			zoneRepo := zone.RepositoryMock{}
			relayMang := execution.RelayManagerMock{}
			logRepo := execution.LogRepositoryMock{}
			exec := execution.NewExecutor(&zoneRepo, &relayMang, &logRepo, &sender)
			execInTime := execution.NewExecutorInTime(&repo, exec, st, &sender)

			repo.GetExecutionsFunc = func() (*execution.Execution, error) {
				return tt.execut, tt.repositoryErr
			}
			sender.SendFunc = func(message string) error {
				return tt.sendErr
			}
			zoneRepo.FindFunc = func(id string) *zone.Zone {
				return tt.zon
			}
			relayMang.ActivatePinsFunc = func(pins []string) error {
				return nil
			}
			relayMang.DeactivatePinsFunc = func(pins []string) error {
				return nil
			}
			logRepo.SaveFunc = func(l execution.Log) error {
				return nil
			}

			err := execInTime.Execute(tt.t)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
