package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
)

func TestCreateZoneHandle(t *testing.T) {
	err := errors.New("")
	cmd := app.CreateZoneCmd{
		ID:       "id",
		ZoneName: "name",
		Relays:   []string{"1"},
	}
	invalidCmd := app.CreateZoneCmd{
		ID:       "",
		ZoneName: "",
		Relays:   nil,
	}
	zon := fixtures.ZoneBuilder{}.Build()
	tests := []struct {
		name string
		cmd  cqs.Command
		expectedErr, relayErr,
		zoneErr, saveErr, updateErr error
		zone zone.Zone
	}{
		{
			name:        "with an invalid command, then it returns and invalid command error",
			cmd:         &invalidCommand{},
			expectedErr: cqs.InvalidCommandError{},
		},
		{
			name:        "with a valid command and find relay returns an error, then it returns same error",
			cmd:         cmd,
			expectedErr: err,
			relayErr:    err,
		},
		{
			name:        "with a valid command and find relay returns a not found, then it returns a create zone error",
			cmd:         cmd,
			expectedErr: app.CreateZoneError{},
			relayErr:    vo.NotFoundError{},
		},
		{
			name:        "with a valid command and find zone returns an error, then it returns same error",
			cmd:         cmd,
			expectedErr: err,
			zoneErr:     err,
		},
		{
			name:        "with a valid command and new zone returns an error, then it returns a create zone error",
			cmd:         invalidCmd,
			expectedErr: app.CreateZoneError{},
			zoneErr:     vo.NotFoundError{},
		},
		{
			name:        "with a valid command and save zone returns an error, then it returns same error",
			cmd:         cmd,
			expectedErr: err,
			zoneErr:     vo.NotFoundError{},
			saveErr:     err,
		},
		{
			name:        "with a valid command and update zone returns an error, then it returns same error",
			cmd:         cmd,
			expectedErr: err,
			updateErr:   err,
			zone:        zon,
		},
		{
			name: "with a valid command, then it returns nil",
			cmd:  cmd,
			zone: zon,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a CreateZone command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			rr := &RelayRepositoryMock{
				FindByKeyFunc: func(ctx context.Context, key string) (zone.Relay, error) {
					return zone.Relay{}, tt.relayErr
				},
			}
			zr := &ZoneRepositoryMock{
				FindByIDFunc: func(ctx context.Context, id string) (zone.Zone, error) {
					return tt.zone, tt.zoneErr
				},
				SaveFunc: func(ctx context.Context, zo zone.Zone) error {
					return tt.saveErr
				},
				UpdateFunc: func(ctx context.Context, zo zone.Zone) error {
					return tt.updateErr
				},
			}
			handler := app.NewCreateZone(rr, zr)
			_, errHand := handler.Handle(context.Background(), tt.cmd)
			test.CheckErrorsType(t, tt.expectedErr, errHand)
		})
	}
}

type invalidCommand struct{}

func (i invalidCommand) Name() string {
	return "invalidCommand"
}
