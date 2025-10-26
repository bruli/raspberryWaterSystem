package config

import (
	"fmt"
	"net/url"
	"strconv"
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
	TelegramToken           = ProjectPrefix + "TELEGRAM_TOKEN"
	TelegramChatID          = ProjectPrefix + "TELEGRAM_CHAT_ID"
	TelegramBotEnabled      = ProjectPrefix + "TELEGRAM_BOT_ENABLED"
)

type Config struct {
	serverURL     string
	environment   EnvironmentType
	zonesFile     string
	authToken     string
	rainServerURL url.URL
	dailyProgramsFile, oddProgramsFile,
	evenProgramsFile, weeklyProgramsFile,
	temperatureProgramsFile string
	executionLogsFile  string
	telegramToken      string
	telegramChatID     int
	telegramBotEnabled bool
}

func (c Config) TelegramBotEnabled() bool {
	return c.telegramBotEnabled
}

func (c Config) RainServerURL() url.URL {
	return c.rainServerURL
}

func (c Config) TelegramToken() string {
	return c.telegramToken
}

func (c Config) TelegramChatID() int {
	return c.telegramChatID
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

func (c Config) Environment() EnvironmentType {
	return c.environment
}

func NewConfig() (*Config, error) {
	servUrl, err := Value(ServerURL)
	if err != nil {
		return nil, err
	}
	environ, err := environment()
	if err != nil {
		return nil, err
	}
	zones, err := Value(ZonesFile)
	if err != nil {
		return nil, err
	}
	auth, err := Value(AuthToken)
	if err != nil {
		return nil, err
	}
	rain, err := Value(RainServerURL)
	if err != nil {
		return nil, err
	}
	rainUrl, err := url.Parse(rain)
	if err != nil {
		return nil, nil
	}
	daily, odd, even, weekly, temp, err := programsFiles()
	if err != nil {
		return nil, err
	}
	execLogs, err := Value(ExecutionLogsFile)
	if err != nil {
		return nil, err
	}
	telegramToken, telegramChatID, err := telegram()
	if err != nil {
		return nil, err
	}
	botEnabledStr, err := Value(TelegramBotEnabled)
	if err != nil {
		return nil, err
	}
	botEnabled, err := strconv.ParseBool(botEnabledStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse telegram bot enabled: %s", botEnabledStr)
	}
	return &Config{
		serverURL:               servUrl,
		environment:             environ,
		zonesFile:               zones,
		authToken:               auth,
		rainServerURL:           *rainUrl,
		dailyProgramsFile:       daily,
		oddProgramsFile:         odd,
		evenProgramsFile:        even,
		weeklyProgramsFile:      weekly,
		temperatureProgramsFile: temp,
		executionLogsFile:       execLogs,
		telegramToken:           telegramToken,
		telegramChatID:          telegramChatID,
		telegramBotEnabled:      botEnabled,
	}, nil
}

func environment() (EnvironmentType, error) {
	envStr, err := Value(Environment)
	if err != nil {
		return 0, err
	}
	environ, err := ParseEnvironment(envStr)
	if err != nil {
		return 0, err
	}
	return environ, nil
}

func telegram() (string, int, error) {
	token, err := Value(TelegramToken)
	if err != nil {
		return "", 0, err
	}
	chatIDStr, err := Value(TelegramChatID)
	if err != nil {
		return "", 0, err
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return "", 0, err
	}
	return token, chatID, nil
}

func programsFiles() (string, string, string, string, string, error) {
	daily, err := Value(DailyProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	odd, err := Value(OddProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	even, err := Value(EvenProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	weekly, err := Value(WeeklyProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	temp, err := Value(TemperatureProgramsFile)
	if err != nil {
		return "", "", "", "", "", err
	}
	return daily, odd, even, weekly, temp, nil
}
