package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
)

func RemoveTemperatureProgram(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := strconv.ParseFloat(chi.URLParam(r, "temperature"), 32)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: "invalid temperature value, must be a number",
			})
			return
		}
		if _, err = ch.Handle(r.Context(), app.RemoveTemperatureProgramCommand{Temperature: float32(temp)}); err != nil {
			switch {
			case errors.As(err, &vo.NotFoundError{}):
				WriteErrorResponse(w, http.StatusNotFound)
			default:
				WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		WriteResponse(w, http.StatusOK, nil)
	}
}
