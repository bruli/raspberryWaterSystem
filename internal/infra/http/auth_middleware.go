package http

import (
	"net/http"
)

const AuthorizationHeader = "Authorization"

func AuthMiddleware(authToken string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(AuthorizationHeader)
			if token != authToken {
				WriteErrorResponse(w, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
