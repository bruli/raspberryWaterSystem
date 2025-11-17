package memory_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/fixtures"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
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
			st := fixtures.StatusBuilder{Active: true}.Build()
			err := repo.Save(ctx, st)
			require.NoError(t, err)
			current = &st
			require.True(t, st.IsActive())
		})
		t.Run(`when Find method is called,
		then it returns the current status`, func(t *testing.T) {
			currStatus, err := repo.Find(ctx)
			require.NoError(t, err)
			require.Equal(t, *current, currStatus)
			require.True(t, currStatus.IsActive())
		})
		t.Run(`when update method is called,
		then it update current status`, func(t *testing.T) {
			updated := *current
			light := fixtures.LightBuilder{}.Build()
			updated.Update(fixtures.WeatherBuilder{Raining: true}.Build(), light)
			updated.Deactivate()
			err := repo.Update(ctx, updated)
			require.NoError(t, err)
			currUpdated, err := repo.Find(ctx)
			require.NoError(t, err)
			require.Equal(t, updated, currUpdated)
			require.False(t, currUpdated.IsActive())
		})
	})
}
