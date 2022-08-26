//go:build functional
// +build functional

package functional_test

import (
	"fmt"
	http2 "net/http"
	"testing"

	"github.com/bruli/raspberryWaterSystem/fixtures"
	"github.com/bruli/raspberryWaterSystem/internal/infra/http"
	"github.com/stretchr/testify/require"
)

func runZones(t *testing.T) {
	t.Run(`Given a Create zone endpoint`, func(t *testing.T) {
		t.Run(`when a request is sent without authorization,
		then it returns an unauthorized`, func(t *testing.T) {
			zo := fixtures.ZoneBuilder{}.Build()
			relays := make([]int, len(zo.Relays()))
			for i, r := range zo.Relays() {
				relays[i] = r.Id().Int()
			}
			req := http.CreateZoneRequestJson{
				Id:     zo.Id(),
				Name:   zo.Name(),
				Relays: relays,
			}
			resp, err := buildRequestAndSend(ctx, req, nil, http2.MethodPut, "/zones", cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`when a request is sent with authorization,
		then it returns an ok`, func(t *testing.T) {
			zo := fixtures.ZoneBuilder{}.Build()
			relays := make([]int, len(zo.Relays()))
			for i, r := range zo.Relays() {
				relays[i] = r.Id().Int()
			}
			req := http.CreateZoneRequestJson{
				Id:     zo.Id(),
				Name:   zo.Name(),
				Relays: relays,
			}
			resp, err := buildRequestAndSend(ctx, req, authorizationHeader(), http2.MethodPut, "/zones", cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusOK, resp.StatusCode)
			savedZone = &zo
		})
	})
	t.Run(`Given an execute zone endpoint,
	when a request is sent`, func(t *testing.T) {
		t.Run(`without authorization,
		then it returns an unauthorized`, func(t *testing.T) {
			url := fmt.Sprintf("/zones/%s/execute", savedZone.Id())
			resp, err := buildRequestAndSend(ctx, nil, nil, http2.MethodPost, url, cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`with authorization,
		then it returns an ok`, func(t *testing.T) {
			url := fmt.Sprintf("/zones/%s/execute", savedZone.Id())
			req := http.ExecuteZoneRequestJson{Seconds: 5}
			resp, err := buildRequestAndSend(ctx, req, authorizationHeader(), http2.MethodPost, url, cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusOK, resp.StatusCode)
		})
	})
	t.Run(`Given a Remove zone endpoint,
	when a request is sent `, func(t *testing.T) {
		t.Run(`without authentication,
		then it return unauthorized`, func(t *testing.T) {
			resp, err := buildRequestAndSend(ctx, nil, nil, http2.MethodDelete, "/zones/bbf", cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusUnauthorized, resp.StatusCode)
		})
		t.Run(`with authentication,
		then it returns ok`, func(t *testing.T) {
			zo := fixtures.ZoneBuilder{}.Build()
			relays := make([]int, len(zo.Relays()))
			for i, r := range zo.Relays() {
				relays[i] = r.Id().Int()
			}
			req := http.CreateZoneRequestJson{
				Id:     zo.Id(),
				Name:   zo.Name(),
				Relays: relays,
			}
			resp, err := buildRequestAndSend(ctx, req, authorizationHeader(), http2.MethodPut, "/zones", cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusOK, resp.StatusCode)
			resp, err = buildRequestAndSend(ctx, nil, authorizationHeader(), http2.MethodDelete, "/zones/"+zo.Id(), cl)
			require.NoError(t, err)
			require.Equal(t, http2.StatusOK, resp.StatusCode)
		})
	})
}
