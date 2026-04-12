package memory

import (
	"context"
	"sync"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

type StatusRepository struct {
	currentStatus *status.Status
	m             sync.RWMutex
	tracer        trace.Tracer
}

func (s *StatusRepository) Find(ctx context.Context) (status.Status, error) {
	select {
	case <-ctx.Done():
		return status.Status{}, ctx.Err()
	default:
		_, span := s.tracer.Start(ctx, "StatusRepository.Find")
		defer span.End()
		s.m.RLock()
		defer s.m.RUnlock()
		if s.currentStatus == nil {
			err := vo.NotFoundError{}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return status.Status{}, err
		}
		st := s.currentStatus
		span.SetStatus(codes.Ok, "status found")
		return *st, nil
	}
}

func (s *StatusRepository) Update(ctx context.Context, st status.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := s.tracer.Start(ctx, "StatusRepository.Update")
		defer span.End()
		if err := s.Save(ctx, st); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetStatus(codes.Ok, "status updated")
		return nil
	}
}

func (s *StatusRepository) Save(ctx context.Context, st status.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := s.tracer.Start(ctx, "StatusRepository.Save")
		defer span.End()
		s.m.Lock()
		defer s.m.Unlock()
		s.currentStatus = &st
		span.SetStatus(codes.Ok, "status saved")
		return nil
	}
}

func NewStatusRepository(tracer trace.Tracer) *StatusRepository {
	return &StatusRepository{tracer: tracer}
}
