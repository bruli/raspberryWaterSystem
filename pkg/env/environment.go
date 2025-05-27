package env

import "errors"

const (
	DevelopmentEnvironment Environment = iota + 1
	ProductionEnvironment
)

var ErrInvalidEnvironment = errors.New("invalid environment")

type Environment int

func (e Environment) IsProduction() bool {
	return e == ProductionEnvironment
}

var environmentFromStringMap = map[string]Environment{
	"development": DevelopmentEnvironment,
	"production":  ProductionEnvironment,
}

func ParseEnvironment(s string) (Environment, error) {
	e, ok := environmentFromStringMap[s]
	if !ok {
		return 0, ErrInvalidEnvironment
	}
	return e, nil
}
