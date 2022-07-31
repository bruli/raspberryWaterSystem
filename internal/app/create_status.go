package app

import (
	"context"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

const CreateStatusCmdName = "createStatus"

type CreateStatusCmd struct {
	StartedAt time.Time
	Weather   weather.Weather
}

func (c CreateStatusCmd) Name() string {
	return CreateStatusCmdName
}

type CreateStatus struct {
	sr StatusRepository
}

func NewCreateStatus(sr StatusRepository) CreateStatus {
	return CreateStatus{sr: sr}
}

func (c CreateStatus) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, _ := cmd.(CreateStatusCmd)
	st := status.New(co.StartedAt, co.Weather)
	return nil, c.sr.Save(ctx, st)
}
