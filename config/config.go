package config

import "github.com/bruli/raspberryRainSensor/pkg/common/env"

const (
	ProjectPrefix = "WS_"
	ServerURL     = ProjectPrefix + "SERVER_URL"
	Environment   = ProjectPrefix + "ENVIRONMENT"
)

type Config struct {
	serverURL   string
	environment env.Environment
}

func (c Config) ServerURL() string {
	return c.serverURL
}

func (c Config) Environment() env.Environment {
	return c.environment
}

func NewConfig() (Config, error) {
	url, err := env.Value(ServerURL)
	if err != nil {
		return Config{}, err
	}
	envStr, err := env.Value(Environment)
	if err != nil {
		return Config{}, err
	}
	environment, err := env.ParseEnvironment(envStr)
	if err != nil {
		return Config{}, err
	}
	return Config{
		serverURL:   url,
		environment: environment,
	}, nil
}
