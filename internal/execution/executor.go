package execution

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"time"
)

type InvalidExecutorData struct {
	error string
}

func NewInvalidExecutorData(error string) InvalidExecutorData {
	return InvalidExecutorData{error: error}
}

func (i InvalidExecutorData) Error() string {
	return i.error
}

type Executor struct {
	zoneRepository     zone.Repository
	relayManager       RelayManager
	logRepository      LogRepository
	notificationSender NotificationSender
}

func NewExecutor(zoneRepository zone.Repository,
	relayManager RelayManager,
	logRepository LogRepository,
	notificationSender NotificationSender) *Executor {
	return &Executor{zoneRepository: zoneRepository,
		relayManager:       relayManager,
		logRepository:      logRepository,
		notificationSender: notificationSender}
}

func (e *Executor) Execute(seconds uint8, zoneID string) error {
	if seconds == 0 {
		return NewInvalidExecutorData("invalid seconds to execute")
	}

	if zoneID == "" {
		return NewInvalidExecutorData("invalid zone to execute")
	}
	zon := e.zoneRepository.Find(zoneID)
	if zon == nil {
		return fmt.Errorf("'%s' is not a valid zone id", zoneID)
	}

	if err := e.relayManager.ActivatePins(zon.Relays()); err != nil {
		return fmt.Errorf("failed activating pins: %w", err)
	}
	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
	if err := e.relayManager.DeactivatePins(zon.Relays()); err != nil {
		return fmt.Errorf("failed deactivating pins: %w", err)
	}
	l := NewLog(seconds, zoneID, time.Now())
	if err := e.logRepository.Save(*l); err != nil {
		return fmt.Errorf("failed saving execution log: %w", err)
	}

	message := fmt.Sprintf("%s zone executed during %s", zon.Name(), duration.String())
	if err := e.notificationSender.Send(message); err != nil {
		return fmt.Errorf("failed sending execution notification %w", err)
	}
	return nil
}
