package weather

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	"github.com/bruli/raspberryWaterSystem/internal/status"
)

type StatusSetter struct {
	st   *status.Status
	repo Repository
	rain *rain.Reader
}

func NewStatusSetter(st *status.Status, repo Repository, rain *rain.Reader) *StatusSetter {
	return &StatusSetter{st: st, repo: repo, rain: rain}
}

func (s *StatusSetter) Set() error {
	temp, hum, err := reader(s.repo)
	if err != nil {
		return err
	}
	r, err := s.rain.Read()
	if err != nil {
		return fmt.Errorf("failed to read rain: %w", err)
	}
	s.st.SetTemperature(temp)
	s.st.SetHumidity(hum)
	s.st.SetRain(r.IsRain(), r.Value())

	return nil
}
