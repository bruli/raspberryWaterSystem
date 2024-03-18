package memory

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

var currentStatus *status.Status

type StatusRepository struct{}

func (s StatusRepository) Find(ctx context.Context) (status.Status, error) {
	select {
	case <-ctx.Done():
		return status.Status{}, ctx.Err()
	default:
		if currentStatus == nil {
			return status.Status{}, vo.NotFoundError{}
		}
		return *currentStatus, nil
	}
}

func (s StatusRepository) Update(ctx context.Context, st status.Status) error {
	return s.Save(ctx, st)
}

func (s StatusRepository) Save(ctx context.Context, st status.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		currentStatus = &st
		return nil
	}
}
