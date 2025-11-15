package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Metrics() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
