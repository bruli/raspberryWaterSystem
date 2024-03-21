package app

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

type QueryBus struct {
	m map[string]cqs.QueryHandler
}

func (c QueryBus) Handle(ctx context.Context, query cqs.Query) (any, error) {
	hand, ok := c.m[query.Name()]
	if !ok {
		return nil, UnSubscribedQueryError{name: query.Name()}
	}
	return hand.Handle(ctx, query)
}

func NewQueryBus() QueryBus {
	m := make(map[string]cqs.QueryHandler)
	return QueryBus{m: m}
}

func (c QueryBus) Subscribe(name string, query cqs.QueryHandler) {
	c.m[name] = query
}

type UnSubscribedQueryError struct {
	name string
}

func (u UnSubscribedQueryError) Error() string {
	return fmt.Sprintf("query %q not subscribed", u.name)
}
