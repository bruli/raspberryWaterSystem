package memory

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

type StatusRepository struct {
	currentStatus *status.Status
}

func NewStatusRepository() *StatusRepository {
	return &StatusRepository{}
}

func (s *StatusRepository) Find(ctx context.Context) (status.Status, error) {
	select {
	case <-ctx.Done():
		return status.Status{}, ctx.Err()
	default:
		if s.currentStatus == nil {
			return status.Status{}, vo.NotFoundError{}
		}
		st := s.currentStatus
		return *st, nil
	}
}

func (s *StatusRepository) Update(ctx context.Context, st status.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return s.Save(ctx, st)
	}
}

func (s *StatusRepository) Save(ctx context.Context, st status.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.currentStatus = &st
		return nil
	}
}
