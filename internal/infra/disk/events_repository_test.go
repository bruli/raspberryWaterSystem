//go:build infra

package disk_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/infra/disk"
	"github.com/stretchr/testify/require"
)

func TestEventsRepository(t *testing.T) {
	t.Run(`Given an Events repository`, func(t *testing.T) {
		repo := disk.NewEventsRepository("/tmp/events")
		number := 3
		var events []disk.Log
		t.Run(`when Save method is called
		then it create files`, func(t *testing.T) {
			for i := 0; i < number; i++ {
				now := time.Now()
				err := repo.Save(t.Context(), fmt.Sprintf("zone_%v", i), 10, now.Add(time.Minute*time.Duration(i)*10))
				require.NoError(t, err)
			}
		})
		t.Run(`when FindAll method it's called,
		then it return a valid slice`, func(t *testing.T) {
			got, err := repo.FindAll(t.Context())
			require.NoError(t, err)
			require.Len(t, got, number)
			events = got
		})
		t.Run(`when Remove method it's called,
			then it delete the json file`, func(t *testing.T) {
			for i := range events {
				err := repo.Remove(t.Context(), &events[i])
				require.NoError(t, err)
			}
		})
	})
}
