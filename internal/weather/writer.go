package weather

import "fmt"

type Writer struct {
	repo      Repository
	writeRepo WriteRepository
}

func NewWriter(repo Repository, writeRepo WriteRepository) *Writer {
	return &Writer{repo: repo, writeRepo: writeRepo}
}

func (w *Writer) Write() error {
	temp, hum, err := w.repo.Read()
	if err != nil {
		return fmt.Errorf("failed to read weather data: %w", err)
	}

	if err := w.writeRepo.Write(temp, hum); err != nil {
		return fmt.Errorf("failed to write weather data: %w", err)
	}
	return nil
}
