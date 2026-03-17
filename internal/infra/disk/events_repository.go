package disk

import (
	"context"
	"fmt"
	"os"
	"time"
)

type EventsRepository struct {
	eventsDir string
}

func (e EventsRepository) Save(ctx context.Context, zone string, seconds int, executedAt time.Time) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := os.MkdirAll(e.eventsDir, 0755); err != nil {
			return err
		}
		lo := Log{
			Seconds:    seconds,
			ZoneName:   zone,
			ExecutedAt: executedAt,
		}
		return writeJsonFile(fmt.Sprintf("%s/%v.json", e.eventsDir, executedAt.Unix()), lo)
	}
}

func (e EventsRepository) Remove(ctx context.Context, lo *Log) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		unix := lo.ExecutedAt.Unix()
		err := os.Remove(fmt.Sprintf("%s/%v.json", e.eventsDir, unix))
		if err != nil {
			return fmt.Errorf("failed to remove event file: %w", err)
		}
		return nil
	}
}

func (e EventsRepository) FindAll(ctx context.Context) ([]Log, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		files, err := os.ReadDir(e.eventsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read events directory: %w", err)
		}
		logs := make([]Log, len(files))
		for i, f := range files {
			if err := readJsonFile(fmt.Sprintf("%s/%s", e.eventsDir, f.Name()), &logs[i]); err != nil {
				return nil, fmt.Errorf("failed to read event file: %w", err)
			}
		}
		return logs, nil
	}
}

func NewEventsRepository(eventsDir string) *EventsRepository {
	return &EventsRepository{eventsDir: eventsDir}
}
