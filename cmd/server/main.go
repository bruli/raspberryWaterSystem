package main

import (
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	"log"
	"os"
	"strconv"
)

func main() {
	serverURL := os.Getenv("SERVER_URL")
	zonesFile := os.Getenv("ZONES_FILE")
	authToken := os.Getenv("AUTH_TOKEN")
	executionsFile := os.Getenv("EXECUTIONS_FILE")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	devMode := os.Getenv("DEV_MODE")
	rainSensorServerUrl := os.Getenv("RAIN_SENSOR_SERVER_URL")
	chatID, err := strconv.ParseInt(telegramChatID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	dev, err := strconv.ParseBool(devMode)
	if err != nil {
		log.Fatal(err)
	}

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
		rainSensorServerUrl,
		chatID,
		dev,
	)
	s := server.NewServer(conf)
	s.Run()
}
