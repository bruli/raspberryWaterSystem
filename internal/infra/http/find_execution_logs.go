package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

func FindExecutionLogs(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		limitStr := r.URL.Query().Get("limit")
		if len(limitStr) > 0 {
			limitValue, err := strconv.Atoi(limitStr)
			if err != nil {
				WriteErrorResponse(w, http.StatusBadRequest, Error{
					Code:    ErrorCodeInvalidRequest,
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
				WriteErrorResponse(w, http.StatusBadRequest, Error{
					Code:    ErrorCodeInvalidRequest,
					Message: err.Error(),
				})
			default:
				WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		logs, _ := result.([]program.ExecutionLog)
		resp := make([]ExecutionLogItemResponse, len(logs))
		for i, log := range logs {
			resp[i] = ExecutionLogItemResponse{
				ExecutedAt: strconv.Itoa(int(log.ExecutedAt().Unix())),
				Seconds:    log.Seconds().Int(),
				ZoneName:   log.ZoneName(),
			}
		}
		data, _ := json.Marshal(resp)
		WriteResponse(w, http.StatusOK, data)
	}
}
