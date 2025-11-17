package app

import (
	"context"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const UpdateStatusCmdName = "updateStatus"

type UpdateStatusCmd struct {
	Weather weather.Weather
}

func (u UpdateStatusCmd) Name() string {
	return UpdateStatusCmdName
}

type UpdateStatus struct {
	sr StatusRepository
	lr LightRepository
}

func NewUpdateStatus(sr StatusRepository, lr LightRepository) UpdateStatus {
	return UpdateStatus{sr: sr, lr: lr}
}

func (u UpdateStatus) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(UpdateStatusCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(UpdateStatusCmdName, cmd.Name())
	}
	current, err := u.sr.Find(ctx)
	if err != nil {
		return nil, err
	}
	li, err := u.lr.Find(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	current.Update(co.Weather, li)
	return nil, u.sr.Update(ctx, current)
}
