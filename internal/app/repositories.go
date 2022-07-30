package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
)

//go:generate moq -out zmock_repositories_test.go -pkg app_test . ZoneRepository

type ZoneRepository interface {
	FindByID(ctx context.Context, id string) (zone.Zone, error)
	Save(ctx context.Context, zo zone.Zone) error
	Update(ctx context.Context, zo zone.Zone) error
}
