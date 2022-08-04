package app

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const ExecuteZoneCmdName = "executeZone"

type ExecuteZoneCmd struct {
	Seconds uint
	ZoneID  string
}

func (e ExecuteZoneCmd) Name() string {
	return ExecuteZoneCmdName
}

type ExecuteZone struct {
	zr ZoneRepository
}

func (e ExecuteZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(ExecuteZoneCmd)
	zo, err := e.zr.FindByID(ctx, co.ZoneID)
	if err != nil {
		return nil, err
	}
	if err = zo.Execute(co.Seconds); err != nil {
		return nil, ExecuteZoneError{m: err.Error()}
	}
	return zo.Events(), nil
}

func NewExecuteZone(zr ZoneRepository) ExecuteZone {
	return ExecuteZone{zr: zr}
}

type ExecuteZoneError struct {
	m string
}

func (e ExecuteZoneError) Error() string {
	return fmt.Sprintf("failed to execute zone: %s", e.m)
}
