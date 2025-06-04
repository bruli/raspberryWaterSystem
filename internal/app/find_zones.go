package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const FindZonesQueryName = "findZones"

type FindZonesQuery struct{}

func (f FindZonesQuery) Name() string {
	return FindZonesQueryName
}

type FindZones struct {
	zr ZoneRepository
}

func (f FindZones) Handle(ctx context.Context, query cqs.Query) (any, error) {
	return f.zr.FindAll(ctx)
}

func NewFindZones(zr ZoneRepository) *FindZones {
	return &FindZones{zr: zr}
}
