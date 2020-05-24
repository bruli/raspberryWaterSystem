package weather

import "fmt"

func reader(r Repository) (float32, float32, error) {
	temp, hum, err := r.Read()
	if err != nil {
		return 0, 0, fmt.Errorf("failed reading weather data: %w", err)
	}
	return temp, hum, nil
}
