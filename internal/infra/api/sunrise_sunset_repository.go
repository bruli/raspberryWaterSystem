package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/domain/status"
)

const (
	latitude  = "41.545"
	longitude = "2.109"
	urlRaw    = "https://api.sunrise-sunset.org/json?lat=%s&lng=%s&date=%s&formatted=0"
	location  = "Europe/Madrid"
)

type results struct {
	Results struct {
		Sunrise time.Time `json:"sunrise"`
		Sunset  time.Time `json:"sunset"`
	} `json:"results"`
}

type SunriseSunsetRepository struct {
	cl    *http.Client
	cache map[string]*status.Light
	sync.RWMutex
}

func (s *SunriseSunsetRepository) Find(ctx context.Context, date time.Time) (*status.Light, error) {
	day := s.formatDay(date)
	defer func() {
		tomorrow := s.formatDay(date.AddDate(0, 0, 1))
		_, _ = s.getAndCache(ctx, tomorrow)
	}()
	s.RLock()
	light, ok := s.cache[day]
	s.RUnlock()
	if ok {
		return light, nil
	}
	return s.getAndCache(ctx, day)
}

func (s *SunriseSunsetRepository) CleanYesterday(ctx context.Context) {
	tick := time.NewTicker(4 * time.Hour)
	select {
	case <-ctx.Done():
		return
	case <-tick.C:
		yesterday := s.formatDay(time.Now().AddDate(0, 0, -1))
		delete(s.cache, yesterday)
	}
}

func (s *SunriseSunsetRepository) formatDay(date time.Time) string {
	return date.Format("2006-01-02")
}

func (s *SunriseSunsetRepository) getAndCache(ctx context.Context, day string) (*status.Light, error) {
	url := fmt.Sprintf(urlRaw, latitude, longitude, day)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get light: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var re results
	if err = json.Unmarshal(body, &re); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}

	li, err := status.NewLight(re.Results.Sunrise.In(loc), re.Results.Sunset.In(loc))
	if err != nil {
		return nil, fmt.Errorf("failed to parse sunrise light: %w", err)
	}
	s.Lock()
	defer s.Unlock()
	s.cache[day] = li
	return li, nil
}

func NewSunriseSunsetRepository(timeout time.Duration) *SunriseSunsetRepository {
	cl := &http.Client{
		Timeout: timeout,
	}
	return &SunriseSunsetRepository{cl: cl, cache: make(map[string]*status.Light)}
}
