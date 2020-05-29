package server

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/jsontime"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type logBody struct {
	Message   string `json:"message"`
	CreatedAt jsontime.JsonTime  `json:"created_at"`
}

func newLogBody(message string, createdAt jsontime.JsonTime) *logBody {
	return &logBody{Message: message, CreatedAt: createdAt}
}

type logsBody []*logBody

func (b *logsBody) add(l *logBody) {
	*b = append(*b, l)
}

type getExecutionLogs struct {
	response *response
	readLogs *execution.ReadLogs
}

func (g getExecutionLogs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lo, err := g.readLogs.Read()
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}

	resp := logsBody{}
	for _, j := range *lo {
		message := fmt.Sprintf("%s executed during %v", j.Zone, j.Seconds)
		resp.add(newLogBody(message, jsontime.JsonTime(j.CreatedAt)))
	}

	data, err := jsoniter.Marshal(resp)
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}

	g.response.writeJSONResponse(w, http.StatusOK, data)
}

func newGetExecutionLogs(readLogs *execution.ReadLogs, log logger.Logger) *getExecutionLogs {
	return &getExecutionLogs{readLogs: readLogs, response: newResponse(log)}
}
