package app

import (
	"context"
	"errors"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
)

var ErrStatusAlreadyExist = errors.New("status already exist")

const CreateStatusCmdName = "createStatus"

type CreateStatusCmd struct {
	StartedAt vo.Time
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
	_, err := c.sr.Find(ctx)
	if err == nil {
		return nil, ErrStatusAlreadyExist
	}
	switch {
	case errors.As(err, &vo.NotFoundError{}):
		st := status.New(co.StartedAt, co.Weather)
		return nil, c.sr.Save(ctx, st)
	default:
		return nil, err
	}
}
