package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
)

func FindExecutionLogs(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		limitStr := r.URL.Query().Get("limit")
		if len(limitStr) > 0 {
			limitValue, err := strconv.Atoi(limitStr)
			if err != nil {
				httpx.WriteErrorResponse(w, http.StatusBadRequest, httpx.Error{
					Code:    httpx.ErrorCodeInvalidRequest,
					Message: "invalid limit value",
				})
				return
			}
			limit = limitValue
		}
		result, err := qh.Handle(r.Context(), app.FindExecutionLogsQuery{Limit: limit})
		if err != nil {
			switch {
			case errors.Is(err, app.ErrInvalidExecutionsLogLimit):
				httpx.WriteErrorResponse(w, http.StatusBadRequest, httpx.Error{
					Code:    httpx.ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
			default:
				httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		logs, _ := result.([]program.ExecutionLog)
		resp := make([]ExecutionLogItemResponse, len(logs))
		for i, log := range logs {
			resp[i] = ExecutionLogItemResponse{
				ExecutedAt: log.ExecutedAt().EpochString(),
				Seconds:    log.Seconds().Int(),
				ZoneName:   log.ZoneName(),
			}
		}
		data, _ := json.Marshal(resp)
		httpx.WriteResponse(w, http.StatusOK, data)
	}
}
