package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	ExecutionLogsEventName = "execution.logs"
	WeatherEventNam        = "weather"
)

type Event struct {
	ID        string    `json:"id"`
	EventName string    `json:"event_name"`
	EventAt   time.Time `json:"event_at"`
	Payload   []byte    `json:"payload"`
}

func NewFromExecutionLog(lo *Log) (*Event, error) {
	payload, err := json.Marshal(lo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal execution log: %w", err)
	}
	return &Event{
		ID:        fmt.Sprintf("%s-%v", lo.ZoneName, lo.ExecutedAt.Unix()),
		EventName: ExecutionLogsEventName,
		EventAt:   time.Now(),
		Payload:   payload,
	}, nil
}

type EventsRepository struct {
	eventsDir string
}

func (e EventsRepository) Save(ctx context.Context, ev *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := os.MkdirAll(e.eventsDir, 0755); err != nil {
			return err
		}

		return writeJsonFile(fmt.Sprintf("%s/%v.json", e.eventsDir, ev.EventAt.UnixMicro()), ev)
	}
}

func (e EventsRepository) Remove(ctx context.Context, ev *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		unix := ev.EventAt.UnixMicro()
		err := os.Remove(fmt.Sprintf("%s/%v.json", e.eventsDir, unix))
		if err != nil {
			return fmt.Errorf("failed to remove event file: %w", err)
		}
		return nil
	}
}

func (e EventsRepository) FindAll(ctx context.Context) ([]Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		files, err := os.ReadDir(e.eventsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read events directory: %w", err)
		}
		events := make([]Event, len(files))
		for i, f := range files {
			if err := readJsonFile(fmt.Sprintf("%s/%s", e.eventsDir, f.Name()), &events[i]); err != nil {
				return nil, fmt.Errorf("failed to read event file: %w", err)
			}
		}
		return events, nil
	}
}

func NewEventsRepository(eventsDir string) *EventsRepository {
	return &EventsRepository{eventsDir: eventsDir}
}
