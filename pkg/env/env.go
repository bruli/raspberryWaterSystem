package env

import (
	"fmt"
	"os"
)

// EmptyEnvironmentKeyError is self described
type EmptyEnvironmentKeyError struct {
	key string
}

// NewEmptyEnvironmentKeyError is a constructor
func NewEmptyEnvironmentKeyError(key string) EmptyEnvironmentKeyError {
	return EmptyEnvironmentKeyError{key: key}
}

func (i EmptyEnvironmentKeyError) Error() string {
	return fmt.Sprintf("empty value from environment key %q", i.key)
}

// Value read value from environment variable
func Value(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", NewEmptyEnvironmentKeyError(key)
	}
	return value, nil
}
