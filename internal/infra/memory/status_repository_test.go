package memory_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/stretchr/testify/require"

	"github.com/bruli/raspberryWaterSystem/internal/infra/memory"
)

func TestStatusRepository(t *testing.T) {
	ctx := context.Background()
	repo := memory.StatusRepository{}
	var current *status.Status
	_ = current
	t.Run(`Given a status repository,`, func(t *testing.T) {
		t.Run(`when Find method is called,
		then it returns a not found error`, func(t *testing.T) {
			_, err := repo.Find(ctx)
			require.ErrorAs(t, err, &vo.NotFoundError{})
		})
		t.Run(`when Save method is called,
		then it save status and not return any error`, func(t *testing.T) {
			st := fixtures.StatusBuilder{}.Build()
			err := repo.Save(ctx, st)
			require.NoError(t, err)
			current = &st
		})
		t.Run(`when Find method is called,
		then it returns the current status`, func(t *testing.T) {
			currStatus, err := repo.Find(ctx)
			require.NoError(t, err)
			require.Equal(t, *current, currStatus)
		})
		t.Run(`when update method is called,
		then it update current status`, func(t *testing.T) {
			updated := *current
			updated.Update(fixtures.WeatherBuilder{Raining: true}.Build())
			err := repo.Update(ctx, updated)
			require.NoError(t, err)
			currUpdated, err := repo.Find(ctx)
			require.NoError(t, err)
			require.Equal(t, updated, currUpdated)
		})
	})
}
