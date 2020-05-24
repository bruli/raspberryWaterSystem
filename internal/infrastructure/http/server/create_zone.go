package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type ZoneBody struct {
	ID     string
	Name   string
	Relays []string
}

type createZone struct {
	createZone *zone.Creator
	body       *ZoneBody
	response   *response
}

func newCreateZone(create *zone.Creator, logger logger.Logger) *createZone {
	return &createZone{createZone: create, body: &ZoneBody{}, response: newResponse(logger)}
}

func (h *createZone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)

	b := h.body
	if err := decoder.Decode(b); err != nil {
		h.response.generateJSONErrorResponse(w, err)
		return
	}

	if err := h.createZone.Create(b.ID, b.Name, b.Relays); err != nil {
		h.response.generateJSONErrorResponse(w, err)
		return
	}
	h.response.writeJSONResponse(w, http.StatusAccepted, nil)
}
