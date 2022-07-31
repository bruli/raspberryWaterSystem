package memory

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

type StatusRepository struct {
	currentStatus *status.Status
}

func (s StatusRepository) Save(ctx context.Context, st status.Status) error {
	s.currentStatus = &st
	return nil
}
