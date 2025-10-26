package env_test

import (
	"testing"

	"github.com/bruli/raspberryWaterSystem/internal/config"
	"github.com/stretchr/testify/require"
)

func TestParseEnvironment(t *testing.T) {
	tests := []struct {
		name, value  string
		expected     config.EnvironmentType
		isProduction bool
	}{
		{
			name:  "with an invalid value, then it returns an invalid environment error",
			value: "invalid",
		},
		{
			name:         "with a valid development value, then it returns a development environment",
			value:        "development",
			expected:     config.DevelopmentEnvironment,
			isProduction: false,
		},
		{
			name:         "with a valid production value, then it returns a production environment",
			value:        "production",
			expected:     config.ProductionEnvironment,
			isProduction: true,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a ParseEnvironment function,
		when is calle `+tt.name, func(t *testing.T) {
			t.Parallel()
			environment, err := config.ParseEnvironment(tt.value)
			if err != nil {
				require.ErrorIs(t, err, config.ErrInvalidEnvironment)
				return
			}
			require.Equal(t, tt.expected, environment)
			require.Equal(t, tt.isProduction, environment.IsProduction())
		})
	}
}
