package server

import (
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/zone"
	jsoniter "github.com/json-iterator/go"
)

type getZones struct {
	body     *ZonesBody
	getZones *zone.Getter
	response *response
}

func newGetZones(getter *zone.Getter, logger logger.Logger) *getZones {
	return &getZones{
		body:     &ZonesBody{},
		getZones: getter,
		response: newResponse(logger)}
}

func (h *getZones) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zo := h.getZones.Get()

	data, err := jsoniter.Marshal(h.buildZonesBody(zo))
	if err != nil {
		h.response.generateJSONErrorResponse(w, err)
	}

	h.response.writeJSONResponse(w, http.StatusOK, data)
}

func (h *getZones) buildZonesBody(zones *zone.Zones) *ZonesBody {
	for _, j := range *zones {
		h.body.add(&ZoneBody{ID: j.Id(), Name: j.Name(), Relays: j.Relays()})
	}

	return h.body
}

type ZonesBody []*ZoneBody

func (zo *ZonesBody) add(z *ZoneBody) {
	*zo = append(*zo, z)
}
