package acceptance

import (
	"bytes"
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	serverURL      = "http://0.0.0.0:5555"
	zonesFile      = "./assets/zones.yml"
	authToken      = "auth_token"
	executionsFile = "./assets/executions.yml"
	mysqlHost      = "localhost"
	mysqlPort      = "3306"
	mysqlUser      = "raspberry"
	mysqlPass      = "raspberry"
	mysqlDatabase  = "raspberryWaterSystem"
	telegramToken  = "telegram"
	telegramChatID = 12345
	dev            = true
	rainSensorURL  = "http:192.168.1.10"
)

func getConfig() *server.Config {
	return server.NewConfig(serverURL,
		zonesFile,
		authToken,
		executionsFile,
		mysqlHost,
		mysqlPort,
		mysqlUser,
		mysqlPass,
		mysqlDatabase,
		telegramToken,
		rainSensorURL,
		telegramChatID,
		dev,
	)
}

func TestHomepage(t *testing.T) {
	tests := map[string]struct {
		responseCode int
	}{
		"it should return OK status": {responseCode: 200},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := sendRequest("GET", "/", nil)
			assert.Nil(t, err)
			assert.Equal(t, tt.responseCode, resp.StatusCode)
		})
	}
}

func sendRequest(method, endpoint string, body []byte) (*http.Response, error) {
	config := getConfig()
	url := config.ServerURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Add("Authorization", config.AuthToken)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, err
}
