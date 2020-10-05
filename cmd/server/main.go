package main

import (
	"flag"
	"log"

	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	"github.com/spf13/viper"
)

func main() {
	configFile := flag.String("config", "", "config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("invalid config file: %s", err)
	}

	serverURL := viper.GetString("server_url")
	zonesFile := viper.GetString("zones_file")
	authToken := viper.GetString("auth_token")
	executionsFile := viper.GetString("executions_file")
	telegramToken := viper.GetString("telegram_token")
	telegramChatID := viper.GetInt64("telegram_chat_id")
	mysqlHost := viper.GetString("mysql_host")
	mysqlPort := viper.GetString("mysql_port")
	mysqlUser := viper.GetString("mysql_user")
	mysqlPass := viper.GetString("mysql_pass")
	mysqlDatabase := viper.GetString("mysql_database")
	devMode := viper.GetBool("dev_mode")
	rainSensorServerURL := viper.GetString("rain_sensor_server_url")

	conf := server.NewConfig(serverURL,
		zonesFile,
		authToken,
		executionsFile,
		mysqlHost,
		mysqlPort,
		mysqlUser,
		mysqlPass,
		mysqlDatabase,
		telegramToken,
		rainSensorServerURL,
		telegramChatID,
		devMode,
	)
	s := server.NewServer(conf)
	s.Run()
}
