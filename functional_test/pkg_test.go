//go:build functional
// +build functional

package functional_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bruli/raspberryWaterSystem/pkg/ws"
	"github.com/stretchr/testify/require"
)

func runPkg(t *testing.T) {
	t.Run(`Given a Water system library,`, func(t *testing.T) {
		ctx = context.Background()
		url, err := url.Parse("http://localhost:8083")
		require.NoError(t, err)
		pkg := ws.New(*url, &http.Client{Timeout: 3 * time.Second}, "WT7*P6Yn^2-Y*V*C-h&K6*b!@=HCzhd+")
		t.Run(`when GetStatus method is called,
		then it returns a valid status object`, func(t *testing.T) {
			status, err := pkg.GetStatus(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, status)
		})
	})
}
