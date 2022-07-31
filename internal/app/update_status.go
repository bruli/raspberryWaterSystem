package app

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
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
}

func NewUpdateStatus(sr StatusRepository) UpdateStatus {
	return UpdateStatus{sr: sr}
}

func (u UpdateStatus) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(UpdateStatusCmd)
	current, err := u.sr.Find(ctx)
	if err != nil {
		return nil, err
	}
	current.Update(co.Weather)
	return nil, u.sr.Update(ctx, current)
}
