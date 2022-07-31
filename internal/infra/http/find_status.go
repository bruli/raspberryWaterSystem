package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
			str := strconv.Itoa(int(currenSt.UpdatedAt().Unix()))
			updated = &str
		}
		resp := StatusResponseJson{
			Humidity:        float64(currenSt.Weather().Humidity()),
			IsRaining:       currenSt.Weather().IsRaining(),
			SystemStartedAt: strconv.Itoa(int(currenSt.SystemStartedAt().Unix())),
			Temperature:     float64(currenSt.Weather().Temp()),
			UpdatedAt:       updated,
		}
		data, _ := json.Marshal(resp)
		httpx.WriteResponse(w, http.StatusOK, data)
	}
}