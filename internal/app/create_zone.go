package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/davecgh/go-spew/spew"
)

const CreateZoneCmdName = "createZone"

type CreateZoneCmd struct {
	ID, ZoneName string
	Relays       []int
}

func (c CreateZoneCmd) Name() string {
	return CreateZoneCmdName
}

type CreateZone struct {
	zr ZoneRepository
}

func NewCreateZone(zr ZoneRepository) CreateZone {
	return CreateZone{zr: zr}
}

func (c CreateZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateZoneCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateZoneCmdName, cmd.Name())
	}
	_, err := c.zr.FindByID(ctx, co.ID)
	if err == nil {
		return nil, CreateZoneError{msg: fmt.Sprintf("a zone with id %s, already exists", co.ID)}
	}
	switch {
	case errors.As(err, &vo.NotFoundError{}):
		return nil, c.createZone(ctx, co)
	default:
		return nil, err
	}
}

func (c CreateZone) createZone(ctx context.Context, co CreateZoneCmd) error {
	relays := make([]zone.Relay, len(co.Relays))
	for i, re := range co.Relays {
		r, err := zone.ParseRelay(re)
		if err != nil {
			return CreateZoneError{msg: err.Error()}
		}
		relays[i] = r
	}
	zo, err := zone.New(co.ID, co.ZoneName, relays)
	if err != nil {
		return CreateZoneError{msg: err.Error()}
	}
	return c.zr.Save(ctx, zo)
}

type CreateZoneError struct {
	msg string
}

func (c CreateZoneError) Error() string {
	return spew.Sprintf("failed to create zone: %s", c.msg)
}
