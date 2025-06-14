package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const FindStatusQueryName = "findStatus"

type FindStatusQuery struct{}

func (f FindStatusQuery) Name() string {
	return FindStatusQueryName
}

type FindStatus struct {
	sr StatusRepository
}

func (f FindStatus) Handle(ctx context.Context, _ cqs.Query) (any, error) {
	return f.sr.Find(ctx)
}

func NewFindStatus(sr StatusRepository) FindStatus {
	return FindStatus{sr: sr}
}
