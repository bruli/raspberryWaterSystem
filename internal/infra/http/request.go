package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadRequest(w http.ResponseWriter, r *http.Request, req interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code, errs := buildInvalidRequestErrorsResponse(fmt.Sprintf("invalid request: %s", err.Error()))
		WriteErrorResponse(w, code, errs...)
		return err
	}
	return nil
}
