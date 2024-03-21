package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/go-chi/chi/v5"
)

func ExecuteZone(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteZoneRequestJson
		if err := httpx.ReadRequest(w, r, &req); err != nil {
			return
		}
		id := chi.URLParam(r, "id")
		if _, err := ch.Handle(r.Context(), app.ExecuteZoneCmd{
			Seconds: uint(req.Seconds),
			ZoneID:  id,
		}); err != nil {
			switch {
			case errors.As(err, &vo.NotFoundError{}):
				httpx.WriteErrorResponse(w, http.StatusNotFound)
			default:
				httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		httpx.WriteResponse(w, http.StatusOK, nil)
	}
}
