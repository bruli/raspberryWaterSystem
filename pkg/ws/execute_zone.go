package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

var ErrUnknownZoneToExecute = errors.New("unknown zone to execute")

func ExecuteZone(cl client) ExecuteZoneFunc {
	return func(ctx context.Context, zone string, seconds int) error {
		url := fmt.Sprintf("%s/zones/%s/execute", cl.serverURL.String(), zone)
		req := http2.ExecuteZoneRequestJson{
			Seconds: seconds,
		}
		resp, err := buildRequestAndSend(ctx, http.MethodPost, req, url, cl.token, cl.cl)
		if err != nil {
			return ErrServer
		}
		defer func() { _ = resp.Body.Close() }()
		switch resp.StatusCode {
		case http.StatusNotFound:
			return ErrUnknownZoneToExecute
		case http.StatusInternalServerError:
			return ErrRemoteServerErr
		default:
			return nil
		}
	}
}
