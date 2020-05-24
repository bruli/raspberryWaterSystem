package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/jsontime"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type homePage struct {
	status   *status.Status
	response *response
}

func newHomePage(logger logger.Logger, status *status.Status) *homePage {
	return &homePage{status: status, response: &response{logger: logger}}
}

func (h *homePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, _ := jsoniter.Marshal(h.buildResponse())
	h.response.writeJSONResponse(w, http.StatusOK, response)
}

func (h *homePage) buildResponse() *homePageResponse {
	response := newHomePageResponse()
	start := jsontime.JsonTime(h.status.SystemStarted())
	response.SystemStarted = &start
	response.Temperature = h.status.Temperature()
	response.Humidity = h.status.Humidity()
	response.OnWater = h.status.OnWater()
	response.Rain.IsRaining = h.status.Rain().IsRain()
	response.Rain.Value = h.status.Rain().Value()

	return response
}

type rainResponse struct {
	IsRaining bool`json:"is_raining"`
	Value     uint16 `json:"value"`
}
type homePageResponse struct {
	SystemStarted *jsontime.JsonTime `json:"system_started"`
	Temperature   float32 `json:"temperature"`
	Humidity      float32 `json:"humidity"`
	OnWater       bool `json:"on_water"`
	Rain          *rainResponse `json:"rain"`
}

func newHomePageResponse() *homePageResponse {
	return &homePageResponse{Rain: &rainResponse{}}
}
