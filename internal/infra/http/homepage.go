package http

import (
	"net/http"
)

func Homepage() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteResponse(w, http.StatusOK, nil)
	}
}
