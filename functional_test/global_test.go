//go:build functional
// +build functional

package functional_test

import (
	"context"
	"net/http"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
)

const serverURL = "http://localhost:8083"

var (
	ctx       context.Context
	cl        http.Client
	savedZone *zone.Zone
)

func init() {
	ctx = context.Background()
	cl = http.Client{Timeout: 3 * time.Second}
}
