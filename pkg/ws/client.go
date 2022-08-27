package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	http2 "github.com/bruli/raspberryWaterSystem/internal/infra/http"
)

type client struct {
	cl        HTTPClient
	serverURL url.URL
	token     string
}

func buildRequestAndSend(ctx context.Context, method string, reqBody interface{}, serverUrl string, token string, cl HTTPClient) (*http.Response, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	buff.Write(body)
	req, _ := http.NewRequestWithContext(ctx, method, serverUrl, &buff)
	req.Header.Add(http2.AuthorizationHeader, token)
	return cl.Do(req)
}
