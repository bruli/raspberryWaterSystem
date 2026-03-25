package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	ExecutionLogsEventType  = "execution_logs"
	TerraceWeatherEventType = "terrace_weather"
)

type Event struct {
	ID        string            `json:"id"`
	EventType string            `json:"event_type"`
	EventAt   time.Time         `json:"event_at"`
	Payload   []byte            `json:"payload"`
	Trace     map[string]string `json:"trace"`
}

type Weather struct {
	Temperature float32 `json:"temperature"`
	IsRaining   bool    `json:"is_raining"`
}

func NewFromWeather(ctx context.Context, w *Weather) (*Event, error) {
	payload, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal weather: %w", err)
	}
	return &Event{
		ID:        uuid.NewString(),
		EventType: TerraceWeatherEventType,
		EventAt:   time.Now(),
		Payload:   payload,
		Trace:     buildTracingMap(ctx),
	}, nil
}

func NewFromExecutionLog(ctx context.Context, lo *Log) (*Event, error) {
	payload, err := json.Marshal(lo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal execution log: %w", err)
	}
	return &Event{
		ID:        fmt.Sprintf("%s-%v", lo.ZoneName, lo.ExecutedAt.Unix()),
		EventType: ExecutionLogsEventType,
		EventAt:   time.Now(),
		Payload:   payload,
		Trace:     buildTracingMap(ctx),
	}, nil
}

func buildTracingMap(ctx context.Context) map[string]string {
	carrier := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(carrier))
	return carrier
}

type EventsRepository struct {
	eventsDir string
	tracer    trace.Tracer
}

func (e EventsRepository) Save(ctx context.Context, ev *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := e.tracer.Start(ctx, "EventsRepository.Save")
		defer span.End()
		if err := os.MkdirAll(e.eventsDir, 0o755); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		if err := writeJsonFile(fmt.Sprintf("%s/%v.json", e.eventsDir, ev.EventAt.UnixMicro()), ev); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetStatus(codes.Ok, "event saved")
		return nil
	}
}

func (e EventsRepository) Remove(ctx context.Context, ev *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := e.tracer.Start(ctx, "EventsRepository.Remove")
		defer span.End()
		unix := ev.EventAt.UnixMicro()
		err := os.Remove(fmt.Sprintf("%s/%v.json", e.eventsDir, unix))
		if err != nil {
			err = fmt.Errorf("failed to remove event file: %w", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetStatus(codes.Ok, "event removed")
		return nil
	}
}

func (e EventsRepository) FindAll(ctx context.Context) ([]Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		_, span := e.tracer.Start(ctx, "EventsRepository.FindAll")
		defer span.End()
		files, err := os.ReadDir(e.eventsDir)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("failed to read events directory: %w", err)
		}
		events := make([]Event, len(files))
		for i, f := range files {
			if err := readJsonFile(fmt.Sprintf("%s/%s", e.eventsDir, f.Name()), &events[i]); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return nil, fmt.Errorf("failed to read event file: %w", err)
			}
		}
		span.SetStatus(codes.Ok, "events found")
		return events, nil
	}
}

func NewEventsRepository(eventsDir string, tracer trace.Tracer) *EventsRepository {
	return &EventsRepository{eventsDir: eventsDir, tracer: tracer}
}
