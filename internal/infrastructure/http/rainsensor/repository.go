package rainsensor

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type response struct {
	IsRaining bool   `json:"is_raining"`
	Value     uint16 `json:"value"`
}

type Repository struct {
	endpoint string
}

func NewRepository(serverURL string) *Repository {
	return &Repository{endpoint: fmt.Sprintf("%s/rain", serverURL)}
}

func (r *Repository) Get() (rain.Rain, error) {
	req, err := http.NewRequest(http.MethodGet, r.endpoint, nil)
	if err != nil {
		return rain.Rain{}, fmt.Errorf("error reading rain values: %w", err)
	}

	cl := http.DefaultClient
	cl.Timeout = 5 * time.Second
	res, err := cl.Do(req)
	if err != nil {
		return rain.Rain{}, fmt.Errorf("failed request to rain sensor: %w", err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	resp := response{}
	if err := jsoniter.Unmarshal(body, &resp); err != nil {
		return rain.Rain{}, fmt.Errorf("error unmarshalling body: %w", err)
	}
	return rain.New(resp.IsRaining, resp.Value), nil
}
