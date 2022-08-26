package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

const RemoveZoneCmdName = "removeZone"

type RemoveZoneCmd struct {
	ID string
}

func (r RemoveZoneCmd) Name() string {
	return RemoveZoneCmdName
}

type RemoveZone struct {
	zr ZoneRepository
}

func (r RemoveZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(RemoveZoneCmd)
	zo, err := r.zr.FindByID(ctx, co.ID)
	if err != nil {
		return nil, err
	}
	return nil, r.zr.Remove(ctx, zo)
}

func NewRemoveZone(zr ZoneRepository) RemoveZone {
	return RemoveZone{zr: zr}
}
