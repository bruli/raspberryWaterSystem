package api

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/bruli/raspberryRainSensor/pkg/rs"
)

type RainRepository struct {
	handler rs.Handler
}

func NewRainRepository(url url.URL) RainRepository {
	cl := &http.Client{Timeout: 3 * time.Second}
	return RainRepository{handler: rs.New(url.String(), cl)}
}

func (r RainRepository) Find(ctx context.Context) (bool, error) {
	rain, err := r.handler.ReadRain(ctx)
	if err != nil {
		return false, nil
	}
	return rain.IsRaining, nil
}
