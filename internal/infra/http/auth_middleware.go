package http

import (
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
)

func AuthMiddleware(authToken string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != authToken {
				httpx.WriteErrorResponse(w, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
