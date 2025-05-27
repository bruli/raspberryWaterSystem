package http

import (
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/go-chi/chi/v5"
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
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: "invalid action name",
			})
			return
		}
		if _, err := ch.Handle(r.Context(), app.ActivateDeactivateServerCmd{Active: active}); err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		WriteResponse(w, http.StatusOK, nil)
	}
}
