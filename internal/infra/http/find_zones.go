package http

import (
	"encoding/json"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"net/http"
)

func FinZones(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := qh.Handle(r.Context(), app.FindZonesQuery{})
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		zones, _ := result.([]zone.Zone)
		resp := make([]ZonesItemResponse, len(zones))
		for i, zon := range zones {
			relays := make([]int, len(zon.Relays()))
			for n, rel := range zon.Relays() {
				relays[n] = rel.Id().Int()
			}
			resp[i] = ZonesItemResponse{
				Id:     zon.Id(),
				Name:   zon.Name(),
				Relays: relays,
			}
		}
		data, _ := json.Marshal(resp)
		WriteResponse(w, http.StatusOK, data)
	}
}
