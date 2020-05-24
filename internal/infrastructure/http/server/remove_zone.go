package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"github.com/gorilla/mux"
	"net/http"
)

type removeZone struct {
	response *response
	remover  *zone.Remover
}

func (r *removeZone) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	zoneId := params["zone_id"]
	err := r.remover.Remove(zoneId)
	if err != nil {
		r.response.generateJSONErrorResponse(w, err)
		return
	}
	r.response.writeJSONResponse(w, http.StatusAccepted, nil)
}

func newRemoveZone(remover *zone.Remover, log logger.Logger) *removeZone {
	return &removeZone{remover: remover, response: newResponse(log)}
}
