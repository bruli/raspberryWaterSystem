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
	//req, err := http.NewRequest(http.MethodGet, r.url.String(), nil)
	//if err != nil {
	//	return false, fmt.Errorf("error reading rain values: %w", err)
	//}
	//
	//cl := http.DefaultClient
	//cl.Timeout = 5 * time.Second
	//res, err := cl.Do(req)
	//if err != nil {
	//	return false, fmt.Errorf("failed request to rain sensor: %w", err)
	//}
	//body, _ := ioutil.ReadAll(res.Body)
	//
	//resp := rainResponse{}
	//if err := json.Unmarshal(body, &resp); err != nil {
	//	return false, fmt.Errorf("error unmarshalling body: %w", err)
	//}
	//return resp.IsRaining, nil

	rain, err := r.handler.ReadRain(ctx)
	if err != nil {
		return false, err
	}
	return rain.IsRaining, nil
}
