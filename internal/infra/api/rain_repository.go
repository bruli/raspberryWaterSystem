package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type rainResponse struct {
	IsRaining bool   `json:"is_raining"`
	Value     uint16 `json:"value"`
}

type RainRepository struct {
	url url.URL
}

func (r RainRepository) Find(ctx context.Context) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, r.url.String(), nil)
	if err != nil {
		return false, fmt.Errorf("error reading rain values: %w", err)
	}

	cl := http.DefaultClient
	cl.Timeout = 5 * time.Second
	res, err := cl.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed request to rain sensor: %w", err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	resp := rainResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, fmt.Errorf("error unmarshalling body: %w", err)
	}
	return resp.IsRaining, nil
}
