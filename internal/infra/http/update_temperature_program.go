package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
	"github.com/go-chi/chi/v5"
)

func UpdateTemperatureProgram(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := strconv.ParseFloat(chi.URLParam(r, "temperature"), 64)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: "invalid temperature value, must be a number",
			})
			return
		}
		var req UpdateTemperatureProgramRequestJson
		if err = ReadRequest(w, r, &req); err != nil {
			return
		}

		prgms, err := buildUpdateTempPrograms(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, Error{
				Code:    ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return
		}

		if _, err = ch.Handle(r.Context(), app.UpdateTemperatureProgramCommand{
			Temperature: float32(temp),
			Programs:    prgms,
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

func buildUpdateTempPrograms(req UpdateTemperatureProgramRequestJson) ([]program.Program, error) {
	prgms := make([]program.Program, len(req))
	for i, p := range req {
		hour, err := program.ParseHour(p.Hour)
		if err != nil {
			return nil, err
		}
		exec := make([]program.Execution, len(p.Executions))
		for n, ex := range p.Executions {
			zo := make([]string, len(ex.Zones))
			for j, z := range ex.Zones {
				zo[j] = z
			}
			sec, err := program.ParseSeconds(ex.Seconds)
			if err != nil {
				return nil, err
			}
			exe, err := program.NewExecution(sec, zo)
			if err != nil {
				return nil, err
			}
			exec[n] = exe
		}
		pr, err := program.New(hour, exec)
		if err != nil {
			return nil, err
		}
		prgms[i] = *pr
	}
	return prgms, nil
}
