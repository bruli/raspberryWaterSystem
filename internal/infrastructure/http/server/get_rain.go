package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type getRain struct {
	read *rain.Reader
	response *response
	st *status.Status
}

type rainBody struct {
	IsRain 	bool  `json:"is_rain"`
	Value uint16 `json:"value"`
}

func newRainBody(isRain bool, value uint16) *rainBody {
	return &rainBody{IsRain: isRain, Value: value}
}

func newGetRain(read *rain.Reader, log logger.Logger, st *status.Status) *getRain {
	return &getRain{read: read, st: st, response: newResponse(log)}
}

func (g *getRain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ra, err := g.read.Read()
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}
	body := newRainBody(ra.IsRain(), ra.Value())
	data, err := jsoniter.Marshal(body)
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}
	g.st.SetRain(ra.IsRain(), ra.Value())
	g.response.writeJSONResponse(w, http.StatusOK, data)
}

