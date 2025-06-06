package http

import (
	"errors"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func CreateTemperatureProgram(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTemperatureProgramRequestJson
		if err := ReadRequest(w, r, &req); err != nil {
			return
		}
		cmd, err := buildCreateTemperatureProgram(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return
		}
		if _, err := ch.Handle(r.Context(), app.CreateTemperatureProgramCommand{
			Temperature: cmd,
		}); err != nil {
			switch {
			case errors.As(err, &app.CreateTemperatureProgramError{}):
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

func buildCreateTemperatureProgram(req CreateTemperatureProgramRequestJson) (*program.Temperature, error) {
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
	temp, _ := program.NewTemperature(float32(req.Temperature), prgms)
	return temp, nil
}
