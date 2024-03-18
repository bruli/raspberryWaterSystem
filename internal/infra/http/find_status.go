package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bruli/raspberryRainSensor/pkg/common/cqs"
	"github.com/bruli/raspberryRainSensor/pkg/common/httpx"
	"github.com/bruli/raspberryRainSensor/pkg/common/vo"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

func FindStatus(qh cqs.QueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := qh.Handle(r.Context(), app.FindStatusQuery{})
		if err != nil {
			switch {
			case errors.As(err, &vo.NotFoundError{}):
				httpx.WriteErrorResponse(w, http.StatusNotFound)
			default:
				httpx.WriteErrorResponse(w, http.StatusInternalServerError)
			}
			return
		}
		currenSt, _ := result.(status.Status)
		var updated *string
		if currenSt.UpdatedAt() != nil {
			updated = vo.StringPtr(currenSt.UpdatedAt().EpochString())
		}
		resp := StatusResponseJson{
			Active:          currenSt.IsActive(),
			Humidity:        float64(currenSt.Weather().Humidity()),
			IsRaining:       currenSt.Weather().IsRaining(),
			SystemStartedAt: currenSt.SystemStartedAt().EpochString(),
			Temperature:     float64(currenSt.Weather().Temp()),
			UpdatedAt:       updated,
		}
		data, _ := json.Marshal(resp)
		httpx.WriteResponse(w, http.StatusOK, data)
	}
}
