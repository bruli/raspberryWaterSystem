package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
)

func ExecuteZone(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteZoneRequestJson
		if err := ReadRequest(w, r, &req); err != nil {
			return
		}
		id := chi.URLParam(r, "id")
		if _, err := ch.Handle(r.Context(), app.ExecuteZoneCmd{
			Seconds: uint(req.Seconds),
			ZoneID:  id,
		}); err != nil {
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
