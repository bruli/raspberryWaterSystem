package memory

import (
	"context"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"

	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
)

type RelayRepository struct {
	rel map[string]string
}

func NewRelayRepository() RelayRepository {
	return RelayRepository{rel: relays()}
}

func (r RelayRepository) FindByKey(_ context.Context, key string) (zone.Relay, error) {
	rel, ok := r.rel[key]
	if !ok {
		return zone.Relay{}, vo.NewNotFoundError(key)
	}
	return zone.NewRelay(key, rel), nil
}

func relays() map[string]string {
	r := make(map[string]string, 6)
	r["1"] = "18"
	r["2"] = "24"
	r["3"] = "23"
	r["4"] = "25"
	r["5"] = "17"
	r["6"] = "27"

	return r
}
