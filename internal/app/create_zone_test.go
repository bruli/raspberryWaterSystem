package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/test"
	"github.com/bruli/raspberryWaterSystem/internal/app"
)

func TestCreateZoneHandle(t *testing.T) {
	errTest := errors.New("")
	cmd := app.CreateZoneCmd{
		ID:       "id",
		ZoneName: "name",
		Relays:   []int{1},
	}
	tests := []struct {
		name string
		cmd  cqs.Command
		expectedErr, relayErr,
		zoneErr, saveErr, updateErr error
		zone zone.Zone
	}{
		{
			name:        "and find returns nil error, then it returns a create zone error",
			cmd:         cmd,
			expectedErr: app.CreateZoneError{},
		},
		{
			name:        "and find returns an error, then it returns same error",
			cmd:         cmd,
			zoneErr:     errTest,
			expectedErr: errTest,
		},
		{
			name: "with invalid relays, then it returns a create zone error",
			cmd: app.CreateZoneCmd{
				Relays: []int{99},
			},
			zoneErr:     vo.NotFoundError{},
			expectedErr: app.CreateZoneError{},
		},
		{
			name: "with invalid id, then it returns a create zone error",
			cmd: app.CreateZoneCmd{
				Relays: []int{1},
			},
			zoneErr:     vo.NotFoundError{},
			expectedErr: app.CreateZoneError{},
		},
		{
			name:        "and save method returns an error, then it returns same error",
			cmd:         cmd,
			zoneErr:     vo.NotFoundError{},
			saveErr:     errTest,
			expectedErr: errTest,
		},
		{
			name:    "then it returns return nil",
			cmd:     cmd,
			zoneErr: vo.NotFoundError{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(`Given a CreateZone command handler,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			zr := &ZoneRepositoryMock{
				FindByIDFunc: func(ctx context.Context, id string) (zone.Zone, error) {
					return tt.zone, tt.zoneErr
				},
				SaveFunc: func(ctx context.Context, zo zone.Zone) error {
					return tt.saveErr
				},
			}
			handler := app.NewCreateZone(zr)
			_, errHand := handler.Handle(context.Background(), tt.cmd)
			test.CheckErrorsType(t, tt.expectedErr, errHand)
		})
	}
}
