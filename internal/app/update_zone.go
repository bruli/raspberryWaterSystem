package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const UpdateZoneCommandName = "updateZone"

type UpdateZoneCommand struct {
	ID, ZoneName string
	Relays       []int
}

func (u UpdateZoneCommand) Name() string {
	return UpdateZoneCommandName
}

type UpdateZone struct {
	zr ZoneRepository
}

func (u UpdateZone) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(UpdateZoneCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(UpdateZoneCommandName, cmd.Name())
	}
	if _, err := u.zr.FindByID(ctx, co.ID); err != nil {
		switch {
		case errors.As(err, &vo.NotFoundError{}):
			return nil, UpdateZoneError{fmt.Sprintf("a zone with id %s, not found", co.ID)}
		default:
			return nil, err
		}
	}
	relays := make([]zone.Relay, len(co.Relays))
	for i, re := range co.Relays {
		r, err := zone.ParseRelay(re)
		if err != nil {
			return nil, UpdateZoneError{msg: err.Error()}
		}
		relays[i] = r
	}
	zo, err := zone.New(co.ID, co.ZoneName, relays)
	if err != nil {
		return nil, UpdateZoneError{msg: err.Error()}
	}
	return nil, u.zr.Update(ctx, zo)
}

type UpdateZoneError struct {
	msg string
}

func (u UpdateZoneError) Error() string {
	return fmt.Sprintf("failed to update zone: %s", u.msg)
}

func NewUpdateZone(zr ZoneRepository) *UpdateZone {
	return &UpdateZone{zr: zr}
}
