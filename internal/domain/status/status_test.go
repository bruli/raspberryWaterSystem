package status_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run(`Given a Status struct,
	when New function is called, 
	then it returns a valid active status`, func(t *testing.T) {
		start := vo.TimeNow()
		weather := fixtures.WeatherBuilder{}.Build()
		st := status.New(start, weather)
		require.Equal(t, start, st.SystemStartedAt())
		require.Equal(t, weather, st.Weather())
		require.True(t, st.IsActive())

		st.Deactivate()
		require.False(t, st.IsActive())
	})
}
