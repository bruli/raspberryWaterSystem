//go:build infra

package disk_test

import (
	"context"
	"testing"

	"github.com/bruli/raspberryWaterSystem/pkg/vo"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/stretchr/testify/require"
)

func TestExecutionLogRepository(t *testing.T) {
	t.Run(`Given a ExecutionLogRepository,`, func(t *testing.T) {
		ctx := context.Background()
		path := "/tmp/execution_logs.json"
		defer populateFile(t, path)
		repo := disk.NewExecutionLogRepository(path)
		t.Run(`when Save method is called,
		then it not returns error`, func(t *testing.T) {
			logs := []program.ExecutionLog{
				fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 1")}.Build(),
				fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 2")}.Build(),
				fixtures.ExecutionLogBuilder{ZoneName: vo.StringPtr("zone 3")}.Build(),
			}
			err := repo.Save(ctx, logs)
			require.NoError(t, err)
		})
		t.Run(`when FindAll method is called,
		then it returns a valid slice`, func(t *testing.T) {
			logs, err := repo.FindAll(ctx)
			require.NoError(t, err)
			require.Len(t, logs, 3)
		})
	})
}
