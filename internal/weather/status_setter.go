package weather

import (
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
)

type StatusSetter struct {
	st   *status.Status
	repo Repository
	rain *rain.Reader
	log  logger.Logger
}

func NewStatusSetter(st *status.Status, repo Repository, rain *rain.Reader, log logger.Logger) *StatusSetter {
	return &StatusSetter{st: st, repo: repo, rain: rain, log: log}
}

func (s *StatusSetter) Set() error {
	temp, hum, err := reader(s.repo)
	if err != nil {
		return err
	}
	s.st.SetTemperature(temp)
	s.st.SetHumidity(hum)
	r, err := s.rain.Read()
	if err != nil {
		s.log.Fatal(fmt.Errorf("failed to read rain: %w", err))
	}
	s.st.SetRain(r.IsRain(), r.Value())

	return nil
}
