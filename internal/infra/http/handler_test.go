package http_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildRequestBody(body string) *strings.Reader {
	return strings.NewReader(body)
}

//func buildResponse(t *testing.T, body io.ReadCloser, resp interface{}) {
//	respBody, err := ioutil.ReadAll(body)
//	require.NoError(t, err)
//	require.NoError(t, json.Unmarshal(respBody, &resp))
//}

func buildRequestJsonToString(t *testing.T, req interface{}) string {
	d, err := json.Marshal(req)
	require.NoError(t, err)
	return string(d)
}
