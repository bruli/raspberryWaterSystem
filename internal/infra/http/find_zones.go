package http

import (
	"encoding/json"
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/zone"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func FinZones(qh cqs.QueryHandler, tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "FindZonesRequest")
		defer span.End()
		result, err := qh.Handle(ctx, app.FindZonesQuery{})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			WriteErrorResponse(w, http.StatusInternalServerError)
			return
		}
		zones, _ := result.([]*zone.Zone)
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
		span.SetStatus(codes.Ok, "zones found")
		WriteResponse(w, http.StatusOK, data)
	}
}
