package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	DevelopmentEnvironment EnvironmentType = iota + 1
	ProductionEnvironment
)

var ErrInvalidEnvironment = errors.New("invalid environment")

type EmptyEnvironmentKeyError struct {
	key string
}

func NewEmptyEnvironmentKeyError(key string) EmptyEnvironmentKeyError {
	return EmptyEnvironmentKeyError{key: key}
}

func (i EmptyEnvironmentKeyError) Error() string {
	return fmt.Sprintf("empty value from environment key %q", i.key)
}

func Value(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", NewEmptyEnvironmentKeyError(key)
	}
	return value, nil
}

type EnvironmentType int

func (e EnvironmentType) IsProduction() bool {
	return e == ProductionEnvironment
}

var environmentFromStringMap = map[string]EnvironmentType{
	"development": DevelopmentEnvironment,
	"production":  ProductionEnvironment,
}

func ParseEnvironment(s string) (EnvironmentType, error) {
	e, ok := environmentFromStringMap[s]
	if !ok {
		return 0, ErrInvalidEnvironment
	}
	return e, nil
}
