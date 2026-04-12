package cqs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventsRepositoryError is self-described
type EventsRepositoryError struct {
	operation string
	err       error
}

// NewEventsRepositoryError is a constructor
func NewEventsRepositoryError(operation string, err error) EventsRepositoryError {
	return EventsRepositoryError{
		operation: operation,
		err:       err,
	}
}

func (e EventsRepositoryError) Error() string {
	return fmt.Errorf("events repository, %s: %w", e.operation, e.err).Error()
}

// EventFactory must return a typed event.
type EventFactory func() Event

// EventName is self-described
type EventName string

func (e EventName) String() string {
	return string(e)
}

// Event is self-described
//
//go:generate go tool moq -out zmock_event_test.go -pkg cqs_test . Event
type Event interface {
	EventID() uuid.UUID
	EventName() string
	EventAt() time.Time
	AggregateRootID() string
}

var _ Event = BasicEvent{}

// BasicEvent is the minimal domain event struct.
type BasicEvent struct {
	IDAttr              uuid.UUID `json:"id"`
	NameAttr            EventName `json:"name"`
	AtAttr              time.Time `json:"at"`
	AggregateRootIDAttr string    `json:"aggregate_root_id"`
}

// EventID is a getter
func (b BasicEvent) EventID() uuid.UUID {
	return b.IDAttr
}

// EventName is a getter
func (b BasicEvent) EventName() string {
	return b.NameAttr.String()
}

// EventAt is a getter
func (b BasicEvent) EventAt() time.Time {
	return b.AtAttr
}

// AggregateRootID is a getter
func (b BasicEvent) AggregateRootID() string {
	return b.AggregateRootIDAttr
}

// NewBasicEvent is the constructor for the type.
func NewBasicEvent(name EventName, id uuid.UUID, aggRootID string) BasicEvent {
	return BasicEvent{
		IDAttr:              id,
		NameAttr:            name,
		AtAttr:              time.Now(),
		AggregateRootIDAttr: aggRootID,
	}
}

type BasicAggregateRoot struct {
	createdAt time.Time
	events    []Event
}

// CreatedAt is a getter
func (b *BasicAggregateRoot) CreatedAt() time.Time {
	return b.createdAt
}

// NewBasicAggregateRoot is a constructor
func NewBasicAggregateRoot() BasicAggregateRoot {
	return BasicAggregateRoot{
		createdAt: time.Now(),
		events:    nil,
	}
}

func (b *BasicAggregateRoot) Record(evs ...Event) {
	b.events = append(b.events, evs...)
}

// Events is a getter
func (b *BasicAggregateRoot) Events() []Event {
	events := b.events
	b.ClearEvents()
	return events
}

// ClearEvents removes all events from the aggregate root.
// It's exported for testing purposes.
func (b *BasicAggregateRoot) ClearEvents() {
	b.events = nil
}

// Hydrate fills the BasicAggregateRoot fields
func (b *BasicAggregateRoot) Hydrate(createdAt time.Time, events []Event) {
	b.createdAt = createdAt
	b.events = events
}
