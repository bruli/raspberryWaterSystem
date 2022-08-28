package ws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

type Log struct {
	ExecutedAt vo.Time
	Seconds    int
	ZoneName   string
}

func GetLog(cl client) LogsFunc {
	return func(ctx context.Context, number int) ([]Log, error) {
		url := fmt.Sprintf("%s/logs?limit=%v", cl.serverURL.String(), number)
		resp, err := buildRequestAndSend(ctx, http.MethodGet, nil, url, cl.token, cl.cl)
		if err != nil {
			return nil, ErrServer
		}
		defer func() { _ = resp.Body.Close() }()
		switch resp.StatusCode {
		case http.StatusInternalServerError:
			return nil, ErrRemoteServerErr
		case http.StatusBadRequest:
			var errorSchema httpx.ErrorResponseJson
			_ = readResponse(resp, &errorSchema)
			return nil, LogError{e: errorSchema.Errors[0].Message}
		default:
			var logsSchema []http2.ExecutionLogItemResponse
			if err = readResponse(resp, &logsSchema); err != nil {
				return nil, ErrFailedToReadResponse
			}
			logs := make([]Log, len(logsSchema))
			for i, lo := range logsSchema {
				executed, _ := vo.ParseFromEpochStr(lo.ExecutedAt)
				logs[i] = Log{
					ExecutedAt: executed,
					Seconds:    lo.Seconds,
					ZoneName:   lo.ZoneName,
				}
			}
			return logs, nil
		}
	}
}

type LogError struct {
	e string
}

func (l LogError) Error() string {
	return fmt.Sprintf("failed getting logs, error: %s", l.e)
}
