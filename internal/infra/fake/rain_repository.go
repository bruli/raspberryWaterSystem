package fake

import "context"

type RainRepository struct{}

func (r RainRepository) Find(ctx context.Context) (bool, error) {
	return false, nil
}
