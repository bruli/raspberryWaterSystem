package fake

import "context"

type RainRepository struct{}

func (r RainRepository) Find(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return false, nil
	}
}
