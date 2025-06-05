package http

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"net/http"
)

func CreateWeeklyProgram(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateWeeklyProgramRequestJson
		if err := ReadRequest(w, r, &req); err != nil {
			return
		}
		cmd, err := buildCreateWeeklyProgramCommand(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return
		}
		if _, err := ch.Handle(r.Context(), app.CreateWeeklyProgramCommand{Weekly: cmd}); err != nil {
			switch {
			case errors.As(err, &app.CreateWeeklyProgramError{}):
				WriteErrorResponse(w, http.StatusBadRequest, Error{
					Code:    ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
				return
			default:
				WriteErrorResponse(w, http.StatusInternalServerError, Error{})
				return
			}
		}
		WriteResponse(w, http.StatusOK, nil)
	}
}

func buildCreateWeeklyProgramCommand(req CreateWeeklyProgramRequestJson) (*program.Weekly, error) {
	day, err := program.ParseWeekDay(req.WeekDay)
	if err != nil {
		return nil, err
	}
	prgms := make([]program.Program, len(req.Programs))
	for i, p := range req.Programs {
		hour, err := program.ParseHour(p.Hour)
		if err != nil {
			return nil, err
		}
		exec := make([]program.Execution, len(p.Executions))
		for j, e := range p.Executions {
			sec, err := program.ParseSeconds(e.Seconds)
			if err != nil {
				return nil, err
			}
			exe, _ := program.NewExecution(sec, e.Zones)
			exec[j] = exe
		}
		pr, _ := program.New(hour, exec)
		prgms[i] = *pr
	}
	weekly, _ := program.NewWeekly(day, prgms)
	return &weekly, nil
}
