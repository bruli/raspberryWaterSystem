package http

import (
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
)

func Homepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteResponse(w, http.StatusOK, nil)
	}
}
