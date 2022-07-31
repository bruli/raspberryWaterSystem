package config

import (
	"net/url"

	"github.com/bruli/raspberryRainSensor/pkg/common/env"
)

const (
	ProjectPrefix = "WS_"
	ServerURL     = ProjectPrefix + "SERVER_URL"
	Environment   = ProjectPrefix + "ENVIRONMENT"
	ZonesFile     = ProjectPrefix + "ZONES_FILE"
	AuthToken     = ProjectPrefix + "AUTH_TOKEN"
	RainServerURL = ProjectPrefix + "RAIN_SERVER_URL"
)

type Config struct {
	serverURL     string
	environment   env.Environment
	zonesFile     string
	authToken     string
	rainServerURL url.URL
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
	servUrl, err := env.Value(ServerURL)
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
	rain, err := env.Value(RainServerURL)
	if err != nil {
		return Config{}, err
	}
	rainUrl, err := url.Parse(rain)
	if err != nil {
		return Config{}, nil
	}
	return Config{
		serverURL:     servUrl,
		environment:   environment,
		zonesFile:     zones,
		authToken:     auth,
		rainServerURL: *rainUrl,
	}, nil
}
