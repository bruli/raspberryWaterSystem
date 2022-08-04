package http

import (
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

func CreatePrograms(ch cqs.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProgramsRequestJson
		if err := httpx.ReadRequest(w, r, &req); err != nil {
			return
		}
		cmd, err := buildCreateProgramsCmd(req)
		if err != nil {
			httpx.WriteErrorResponse(w, http.StatusBadRequest, httpx.Error{
				Code:    httpx.ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
			return
		}
		if _, err := ch.Handle(r.Context(), cmd); err != nil {
			httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		httpx.WriteResponse(w, http.StatusOK, nil)
	}
}

func buildCreateProgramsCmd(req CreateProgramsRequestJson) (cqs.Command, error) {
	daily, err := buildPrograms(req.Daily)
	if err != nil {
		return nil, err
	}
	odd, err := buildPrograms(req.Odd)
	if err != nil {
		return nil, err
	}
	even, err := buildPrograms(req.Even)
	if err != nil {
		return nil, err
	}
	weekly, err := buildWeeklyPrograms(req.Weekly)
	if err != nil {
		return nil, err
	}
	temp, err := buildTemperaturePrograms(req.Temperature)
	if err != nil {
		return nil, err
	}
	return app.CreateProgramsCmd{
		Daily:       daily,
		Odd:         odd,
		Even:        even,
		Weekly:      weekly,
		Temperature: temp,
	}, nil
}

func buildTemperaturePrograms(requests []TemperatureItemRequest) ([]program.Temperature, error) {
	programs := make([]program.Temperature, len(requests))
	for i, prg := range requests {
		prgms, err := buildPrograms(prg.Programs)
		if err != nil {
			return nil, err
		}
		temp, err := program.NewTemperature(float32(prg.Temperature), prgms)
		if err != nil {
			return nil, err
		}
		programs[i] = temp
	}
	return programs, nil
}

func buildWeeklyPrograms(requests []WeeklyItemRequest) ([]program.Weekly, error) {
	programs := make([]program.Weekly, len(requests))
	for i, prg := range requests {
		day, err := program.ParseWeekDay(prg.WeekDay)
		if err != nil {
			return nil, err
		}
		prgms, err := buildPrograms(prg.Programs)
		if err != nil {
			return nil, err
		}
		weekly, _ := program.NewWeekly(day, prgms)
		programs[i] = weekly
	}
	return programs, nil
}

func buildPrograms(requests []ProgramItemRequest) ([]program.Program, error) {
	programs := make([]program.Program, len(requests))
	for i, prg := range requests {
		hour, err := program.ParseHour(prg.Hour)
		if err != nil {
			return nil, err
		}
		exec := make([]program.Execution, len(prg.Executions))
		for n, ex := range prg.Executions {
			sec, errSec := program.ParseSeconds(ex.Seconds)
			if errSec != nil {
				return nil, errSec
			}
			execution, errExec := program.NewExecution(sec, ex.Zones)
			if errExec != nil {
				return nil, errExec
			}
			exec[n] = execution
		}
		pr, err := program.New(hour, exec)
		if err != nil {
			return nil, err
		}
		programs[i] = pr
	}
	return programs, nil
}
