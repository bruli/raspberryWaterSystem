package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Value string `envconfig:"ENVIRONMENT" required:"true"`
}

type Config struct {
	ServerURL               string `envconfig:"SERVER_URL" required:"true"`
	environment             EnvironmentType
	ZonesFile               string `envconfig:"ZONES_FILE" required:"true"`
	AuthToken               string `envconfig:"AUTH_TOKEN" required:"true"`
	DailyProgramsFile       string `envconfig:"DAILY_PROGRAMS_FILE" required:"true"`
	OddProgramsFile         string `envconfig:"ODD_PROGRAMS_FILE" required:"true"`
	EvenProgramsFile        string `envconfig:"EVEN_PROGRAMS_FILE" required:"true"`
	WeeklyProgramsFile      string `envconfig:"WEEKLY_PROGRAMS_FILE" required:"true"`
	TemperatureProgramsFile string `envconfig:"TEMPERATURE_PROGRAMS_FILE" required:"true"`
	ExecutionLogsFile       string `envconfig:"EXECUTION_LOGS_FILE" required:"true"`
	TelegramToken           string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	TelegramChatID          int    `envconfig:"TELEGRAM_CHAT_ID" required:"true"`
	TelegramBotEnabled      bool   `envconfig:"TELEGRAM_BOT_ENABLED" required:"true"`
}

func New() (*Config, error) {
	var (
		co Config
		e  Environment
	)
	prefix := "WS"

	if err := envconfig.Process(prefix, &co); err != nil {
		return nil, err
	}

	if err := envconfig.Process(prefix, &e); err != nil {
		return nil, err
	}
	env, err := ParseEnvironment(e.Value)
	if err != nil {
		return nil, err
	}
	co.environment = env
	return &co, nil

}
