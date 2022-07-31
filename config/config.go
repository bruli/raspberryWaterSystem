package config

import "github.com/bruli/raspberryRainSensor/pkg/common/env"

const (
	ProjectPrefix = "WS_"
	ServerURL     = ProjectPrefix + "SERVER_URL"
	Environment   = ProjectPrefix + "ENVIRONMENT"
	ZonesFile     = ProjectPrefix + "ZONES_FILE"
	AuthToken     = ProjectPrefix + "AUTH_TOKEN"
)

type Config struct {
	serverURL   string
	environment env.Environment
	zonesFile   string
	authToken   string
}

func (c Config) AuthToken() string {
	return c.authToken
}

func (c Config) ZonesFile() string {
	return c.zonesFile
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
	zones, err := env.Value(ZonesFile)
	if err != nil {
		return Config{}, err
	}
	auth, err := env.Value(AuthToken)
	if err != nil {
		return Config{}, err
	}
	return Config{
		serverURL:   url,
		environment: environment,
		zonesFile:   zones,
		authToken:   auth,
	}, nil
}
