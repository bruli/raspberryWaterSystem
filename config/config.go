package config

import (
	"net/url"

	"github.com/bruli/raspberryRainSensor/pkg/common/env"
)

const (
	ProjectPrefix           = "WS_"
	ServerURL               = ProjectPrefix + "SERVER_URL"
	Environment             = ProjectPrefix + "ENVIRONMENT"
	ZonesFile               = ProjectPrefix + "ZONES_FILE"
	AuthToken               = ProjectPrefix + "AUTH_TOKEN"
	RainServerURL           = ProjectPrefix + "RAIN_SERVER_URL"
	DailyProgramsFile       = ProjectPrefix + "DAILY_PROGRAMS_FILE"
	OddProgramsFile         = ProjectPrefix + "ODD_PROGRAMS_FILE"
	EvenProgramsFile        = ProjectPrefix + "EVEN_PROGRAMS_FILE"
	WeeklyProgramsFile      = ProjectPrefix + "WEEKLY_PROGRAMS_FILE"
	TemperatureProgramsFile = ProjectPrefix + "TEMPERATURE_PROGRAMS_FILE"
	ExecutionLogsFile       = ProjectPrefix + "EXECUTION_LOGS_FILE"
)

type Config struct {
	serverURL     string
	environment   env.Environment
	zonesFile     string
	authToken     string
	rainServerURL url.URL
	dailyProgramsFile, oddProgramsFile,
	evenProgramsFile, weeklyProgramsFile,
	temperatureProgramsFile string
	executionLogsFile string
}

func (c Config) ExecutionLogsFile() string {
	return c.executionLogsFile
}

func (c Config) DailyProgramsFile() string {
	return c.dailyProgramsFile
}

func (c Config) OddProgramsFile() string {
	return c.oddProgramsFile
}

func (c Config) EvenProgramsFile() string {
	return c.evenProgramsFile
}

func (c Config) WeeklyProgramsFile() string {
	return c.weeklyProgramsFile
}

func (c Config) TemperatureProgramsFile() string {
	return c.temperatureProgramsFile
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
	daily, odd, even, weekly, temp, err := programsFiles()
	if err != nil {
		return Config{}, err
	}
	execLogs, err := env.Value(ExecutionLogsFile)
	if err != nil {
		return Config{}, err
	}
	return Config{
		serverURL:               servUrl,
		environment:             environment,
		zonesFile:               zones,
		authToken:               auth,
		rainServerURL:           *rainUrl,
		dailyProgramsFile:       daily,
		oddProgramsFile:         odd,
		evenProgramsFile:        even,
		weeklyProgramsFile:      weekly,
		temperatureProgramsFile: temp,
		executionLogsFile:       execLogs,
	}, nil
}

func programsFiles() (string, string, string, string, string, error) {
	daily, err := env.Value(DailyProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	odd, err := env.Value(OddProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	even, err := env.Value(EvenProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	weekly, err := env.Value(WeeklyProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	temp, err := env.Value(TemperatureProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	return daily, odd, even, weekly, temp, nil
}
