package app

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const ExecuteZoneWithStatusCmdName = "executeZoneWithStatus"

type ExecuteZoneWithStatusCmd struct {
	Seconds uint
	ZoneID  string
}

func (e ExecuteZoneWithStatusCmd) Name() string {
	return ExecuteZoneWithStatusCmdName
}

type ExecuteZoneWithStatus struct {
	zr ZoneRepository
	st StatusRepository
}

func NewExecuteZoneWithStatus(zr ZoneRepository, st StatusRepository) *ExecuteZoneWithStatus {
	return &ExecuteZoneWithStatus{zr: zr, st: st}
}

func (e ExecuteZoneWithStatus) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(ExecuteZoneWithStatusCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(ExecuteZoneWithStatusCmdName, cmd.Name())
	}
	zo, err := e.zr.FindByID(ctx, co.ZoneID)
	if err != nil {
		return nil, err
	}
	st, err := e.st.Find(ctx)
	if err != nil {
		return nil, err
	}
	if err = zo.ExecuteWithStatus(st.IsActive(), st.Weather().IsRaining(), co.Seconds); err != nil {
		return nil, ExecuteZoneWithStatusError{m: err.Error()}
	}
	return zo.Events(), nil
}

type ExecuteZoneWithStatusError struct {
	m string
}

func (e ExecuteZoneWithStatusError) Error() string {
	return fmt.Sprintf("failed to execute zone with status: %s", e.m)
}
