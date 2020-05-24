package server

type Config struct {
	ServerURL, ZonesFile, AuthToken, ExecutionsFile           string
	MysqlHost, MysqlPort, MysqlUser, MysqlPass, MysqlDatabase string
	telegramToken, rainSensorServerUrl                        string
	telegramChatID                                            int64
	devMode                                                   bool
}

func NewConfig(serverURL, zonesFile, authToken, executionsFile string,
	mysqlHost, mysqlPort, mysqlUser, mysqlPass, mysqlDatabase string,
	telegramToken, rainSensorServerUrl string,
	telegramChatID int64,
	devMode bool,
) *Config {
	return &Config{
		ServerURL:           serverURL,
		ZonesFile:           zonesFile,
		AuthToken:           authToken,
		ExecutionsFile:      executionsFile,
		MysqlHost:           mysqlHost,
		MysqlPort:           mysqlPort,
		MysqlUser:           mysqlUser,
		MysqlPass:           mysqlPass,
		MysqlDatabase:       mysqlDatabase,
		telegramChatID:      telegramChatID,
		telegramToken:       telegramToken,
		devMode:             devMode,
		rainSensorServerUrl: rainSensorServerUrl,
	}
}
