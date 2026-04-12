package tracing

import (
	"github.com/bruli/raspberryWaterSystem/internal/cqs"
	"go.opentelemetry.io/otel/trace"
)

type Event struct {
	SpanContext trace.SpanContext
	Event       cqs.Event
}
