package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func CreateZone(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateZoneRequestJson
		if err := ReadRequest(w, r, &req); err != nil {
			return
		}
		if _, err := ch.Handle(r.Context(), app.CreateZoneCmd{
			ID:       req.Id,
			ZoneName: req.Name,
			Relays:   req.Relays,
		}); err != nil {
			switch {
			case errors.As(err, &app.CreateZoneError{}):
				WriteErrorResponse(w, http.StatusBadRequest, Error{
					Code:    ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
			default:
				WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		WriteResponse(w, http.StatusOK, nil)
	}
}
