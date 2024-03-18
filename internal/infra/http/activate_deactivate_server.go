package http

import (
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/go-chi/chi"
)

const (
	ActivateAction   = "activate"
	DeactivateAction = "deactivate"
)

func ActivateDeactivateServer(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := chi.URLParam(r, "action")
		var active bool
		switch action {
		case ActivateAction:
			active = true
		case DeactivateAction:
			active = false
		default:
			httpx.WriteErrorResponse(w, http.StatusBadRequest, httpx.Error{
				Code:    httpx.ErrorCodeInvalidRequest,
				Message: "invalid action name",
			})
			return
		}
		if _, err := ch.Handle(r.Context(), app.ActivateDeactivateServerCmd{Active: active}); err != nil {
			httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		httpx.WriteResponse(w, http.StatusOK, nil)
	}
}
