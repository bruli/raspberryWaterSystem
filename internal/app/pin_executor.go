package app

import "context"

//go:generate moq -out zmock_pin_executor_test.go --pkg app_test . PinExecutor
type PinExecutor interface {
	Execute(ctx context.Context, seconds uint, pins []string) error
}
