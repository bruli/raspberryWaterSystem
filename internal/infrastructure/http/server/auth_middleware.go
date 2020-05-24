package server

import (
	"net/http"
)

type authMiddleware struct {
	authToken string
}

func (a *authMiddleware) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != a.authToken {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newAuthMiddleware(authToken string) *authMiddleware {
	return &authMiddleware{authToken: authToken}
}
