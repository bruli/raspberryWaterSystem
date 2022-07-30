package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
)

func CreateZone(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateZoneRequestJson
		if err := httpx.ReadRequest(w, r, &req); err != nil {
			return
		}
		if _, err := ch.Handle(r.Context(), app.CreateZoneCmd{
			ID:       req.Id,
			ZoneName: req.Name,
			Relays:   req.Relays,
		}); err != nil {
			switch {
			case errors.As(err, &app.CreateZoneError{}):
				httpx.WriteErrorResponse(w, http.StatusBadRequest, httpx.Error{
					Code:    httpx.ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
			default:
				httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		httpx.WriteResponse(w, http.StatusOK, nil)
	}
}
