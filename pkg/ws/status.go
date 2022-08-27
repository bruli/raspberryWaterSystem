package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

type Status struct {
	Humidity        float64
	IsRaining       bool
	SystemStartedAt vo.Time
	Temperature     float64
	UpdatedAt       *vo.Time
}

func GetStatus(cl client) StatusFunc {
	return func(ctx context.Context) (Status, error) {
		url := cl.serverURL.String() + "/status"
		resp, err := buildRequestAndSend(ctx, http.MethodGet, nil, url, cl.token, cl.cl)
		if err != nil {
			return Status{}, ErrServer
		}
		defer func() { _ = resp.Body.Close() }()
		switch resp.StatusCode {
		case http.StatusNotFound:
			return Status{}, ErrRemoteServerErr
		case http.StatusInternalServerError:
			return Status{}, ErrRemoteServerErr
		default:
			var schema http2.StatusResponseJson
			if err = readResponse(resp, &schema); err != nil {
				return Status{}, ErrFailedToReadResponse
			}
			var updated *vo.Time
			if schema.UpdatedAt != nil {
				str, _ := vo.ParseFromEpochStr(vo.StringValue(schema.UpdatedAt))
				updated = vo.TimePtr(str)
			}
			started, _ := vo.ParseFromEpochStr(schema.SystemStartedAt)
			return Status{
				Humidity:        schema.Humidity,
				IsRaining:       schema.IsRaining,
				SystemStartedAt: started,
				Temperature:     schema.Temperature,
				UpdatedAt:       updated,
			}, nil
		}
	}
}

func readResponse(resp *http.Response, schema interface{}) error {
	return json.NewDecoder(resp.Body).Decode(schema)
}
