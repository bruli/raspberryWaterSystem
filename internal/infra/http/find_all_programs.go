package http

import (
	"encoding/json"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"
)

func FindAllPrograms(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := qh.Handle(r.Context(), app.FindAllProgramsQuery{})
		if err != nil {
			httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		programs, _ := result.(app.AllPrograms)
		resp := ProgramsResponseJson{
			Daily:       buildProgramsResponse(programs.Daily),
			Even:        buildProgramsResponse(programs.Even),
			Odd:         buildProgramsResponse(programs.Odd),
			Temperature: buildTemperatureProgramsResponse(programs.Temperature),
			Weekly:      buildWeeklyProgramsResponse(programs.Weekly),
		}
		data, _ := json.Marshal(resp)
		httpx.WriteResponse(w, http.StatusOK, data)
	}
}

func buildWeeklyProgramsResponse(weekly []program.Weekly) []WeeklyItemResponse {
	itemResponse := make([]WeeklyItemResponse, len(weekly))
	for i, we := range weekly {
		itemResponse[i] = WeeklyItemResponse{
			Programs: buildProgramsResponse(we.Programs()),
			WeekDay:  we.WeekDay().String(),
		}
	}
	return itemResponse
}

func buildTemperatureProgramsResponse(temperature []program.Temperature) []TemperatureItemResponse {
	itemResponse := make([]TemperatureItemResponse, len(temperature))
	for i, temp := range temperature {
		itemResponse[i] = TemperatureItemResponse{
			Programs:    buildProgramsResponse(temp.Programs()),
			Temperature: float64(temp.Temperature()),
		}
	}
	return itemResponse
}

func buildProgramsResponse(programs []program.Program) []ProgramItemResponse {
	programItemResponses := make([]ProgramItemResponse, len(programs))
	for i, d := range programs {
		exec := make([]ExecutionItemResponse, len(d.Executions()))
		for n, ex := range d.Executions() {
			zo := make([]string, len(ex.Zones()))
			for x, z := range ex.Zones() {
				zo[x] = z
			}
			exec[n] = ExecutionItemResponse{
				Seconds: ex.Seconds().Int(),
				Zones:   zo,
			}
		}
		programItemResponses[i] = ProgramItemResponse{
			Executions: exec,
			Hour:       d.Hour().String(),
		}
	}
	return programItemResponses
}
