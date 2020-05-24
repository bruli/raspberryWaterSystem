package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

type response struct {
	logger logger.Logger
}

func newResponse(logger logger.Logger) *response {
	return &response{logger: logger}
}

func (r *response) writeJSONErrorResponse(w http.ResponseWriter, code int, err error) {
	message := errorResponse{Error: err.Error()}
	body, _ := jsoniter.Marshal(message)
	r.writeJSONResponse(w, code, body)
}

func (r *response) generateJSONErrorResponse(w http.ResponseWriter, err error) {
	code := r.getCode(err)
	r.logger.Infof("error response. Status code: %v, error body: %w", code, err)
	if http.StatusInternalServerError == code {
		err = errors.New("internal server error")
	}

	r.writeJSONErrorResponse(w, code, err)
}

func (r *response) getCode(err error) int {
	if checkBadRequestErrors(err) {
		return http.StatusBadRequest
	}
	if checkNotFoundErrors(err) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func (r *response) writeJSONResponse(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		r.logger.Fatalf("failed to write response body: %w", err)
	}
}
