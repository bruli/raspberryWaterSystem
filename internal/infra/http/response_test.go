package http_test

import (
	"encoding/json"
	"io/ioutil"
	http2 "net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func readResponse(t *testing.T, resp *http2.Response, schema interface{}) {
	body := resp.Body
	defer func() { _ = body.Close() }()
	respBody, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(respBody, &schema))
}
