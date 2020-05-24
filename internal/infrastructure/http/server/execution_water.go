package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type ExecuteWaterBody struct {
	Seconds uint8  `json:"seconds"`
	Zone    string `json:"zone"`
}

func NewExecuteWaterBody(seconds uint8, zone string) *ExecuteWaterBody {
	return &ExecuteWaterBody{Seconds: seconds, Zone: zone}
}

type ExecutionWater struct {
	data     chan executionData
	response *response
}

func NewExecutionWater(data chan executionData, log logger.Logger) *ExecutionWater {
	return &ExecutionWater{data: data, response: newResponse(log)}
}

func (e *ExecutionWater) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	body := ExecuteWaterBody{}
	if err := decoder.Decode(&body); err != nil {
		e.response.generateJSONErrorResponse(w, err)
		return
	}
	execData := newExecutionData(body.Seconds, body.Zone)
	e.data <- *execData

	e.response.writeJSONResponse(w, http.StatusAccepted, nil)
}
