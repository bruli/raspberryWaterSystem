//go:build functional

package functional

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func buildRequestAndSend(ctx context.Context, requestBody interface{}, headers map[string]string, method, url string, cl http.Client) (*http.Response, error) {
	var buff bytes.Buffer
	if requestBody != nil {
		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body object: %w", err)
		}
		buff.Write(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", serverURL, url), &buff)
	if err != nil {
		return nil, err
	}
	for i, n := range headers {
		req.Header.Set(i, n)
	}
	return cl.Do(req)
}

func authorizationHeader() map[string]string {
	const token = "token"
	return map[string]string{
		"Authorization": token,
	}
}
