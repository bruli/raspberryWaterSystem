package http

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"net/http"
)

func CreateDailyProgram(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProgramRequestJson
		if err := ReadRequest(w, r, &req); err != nil {
			return
		}
		prog, err := buildProgram(w, req)
		if err != nil {
			return
		}
		if _, err = ch.Handle(r.Context(), app.CreateDailyProgramCommand{Program: prog}); err != nil {
			switch {
			case errors.As(err, &app.CreateProgramError{}):
				WriteErrorResponse(w, http.StatusBadRequest, Error{
					Code:    ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
				return
			default:
				WriteErrorResponse(w, http.StatusInternalServerError)
			}
		}
		WriteErrorResponse(w, http.StatusOK)
	}
}

func buildProgram(w http.ResponseWriter, req CreateProgramRequestJson) (*program.Program, error) {
	hour, err := program.ParseHour(req.Hour)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, Error{
			Code:    ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return nil, err
	}
	exec := make([]program.Execution, len(req.Executions))
	for i, ex := range req.Executions {
		sec, err := program.ParseSeconds(ex.Seconds)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return nil, err
		}
		e, err := program.NewExecution(sec, ex.Zones)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
		}
		exec[i] = e
	}
	prog, _ := program.New(hour, exec)
	return prog, err
}
