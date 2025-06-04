package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
)

func RemoveProgram(ch cqs.CommandHandler, prg string) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		hour, err := program.ParseHour(chi.URLParam(r, "hour"))
		if err != nil {
			WriteErrorResponse(writer, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return
		}
		if _, err := ch.Handle(r.Context(), buildRemoveProgramCommand(hour, prg)); err != nil {
			switch {
			case errors.As(err, &vo.NotFoundError{}):
				WriteErrorResponse(writer, http.StatusNotFound)
			default:
				WriteErrorResponse(writer, http.StatusInternalServerError)
			}
		}
		WriteResponse(writer, http.StatusOK, nil)
	}
}

func buildRemoveProgramCommand(hour program.Hour, prg string) cqs.Command {
	switch prg {
	case DailyProgram:
		return app.RemoveDailyProgramCommand{
			Hour: &hour,
		}
	case OddProgram:
		return app.RemoveOddProgramCommand{
			Hour: &hour,
		}
	case EvenProgram:
		return app.RemoveEvenProgramCommand{
			Hour: &hour,
		}
	}
	return nil
}
