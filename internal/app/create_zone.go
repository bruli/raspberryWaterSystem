package app

import (
	"context"
	"errors"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/davecgh/go-spew/spew"
)

const CreateZoneCmdName = "createZone"

type CreateZoneCmd struct {
	ID, ZoneName string
	Relays       []string
}

func (c CreateZoneCmd) Name() string {
	return CreateZoneCmdName
}

type CreateZone struct {
	rr RelayRepository
	zr ZoneRepository
}

func NewCreateZone(rr RelayRepository, zr ZoneRepository) CreateZone {
	return CreateZone{rr: rr, zr: zr}
}

func (c CreateZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateZoneCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateZoneCmdName, cmd.Name())
	}

	if err := c.validateRelays(ctx, co.Relays); err != nil {
		return nil, err
	}
	zo, err := c.zr.FindByID(ctx, co.ID)
	if err != nil {
		if !errors.As(err, &vo.NotFoundError{}) {
			return nil, err
		}
		newZone, err := zone.New(co.ID, co.ZoneName, co.Relays)
		if err != nil {
			return nil, CreateZoneError{msg: err.Error()}
		}
		if err := c.zr.Save(ctx, newZone); err != nil {
			return nil, err
		}
	}
	zo.Update(co.ZoneName, co.Relays)
	if err := c.zr.Update(ctx, zo); err != nil {
		return nil, err
	}
	return nil, nil
}

func (c CreateZone) validateRelays(ctx context.Context, relays []string) error {
	for _, rel := range relays {
		_, err := c.rr.FindByKey(ctx, rel)
		if err != nil {
			if !errors.As(err, &vo.NotFoundError{}) {
				return err
			}
			return CreateZoneError{msg: spew.Sprintf("invalid relay: %s", rel)}
		}
	}
	return nil
}

type CreateZoneError struct {
	msg string
}

func (c CreateZoneError) Error() string {
	return spew.Sprintf("failed to create zone: %s", c.msg)
}
