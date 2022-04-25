package memory_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/memory"
	"github.com/stretchr/testify/require"
)

func TestRelayRepository(t *testing.T) {
	t.Run(`Given a RelayRepository,
	when FindByKey methods is called `, func(t *testing.T) {
		repo := memory.NewRelayRepository()
		ctx := context.Background()
		t.Run(`with a valid key,
		then it returns a relay`, func(t *testing.T) {
			rel, err := repo.FindByKey(ctx, "1")
			require.NoError(t, err)
			require.Equal(t, rel.Pin(), "18")
		})
		t.Run(`with an invalid key,
		then it returns a not found error`, func(t *testing.T) {
			_, err := repo.FindByKey(ctx, "invalid")
			require.ErrorAs(t, err, &vo.NotFoundError{})
		})
	})
}
