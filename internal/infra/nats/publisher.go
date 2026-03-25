package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const StreamName = "WATER_SYSTEM"

type Publisher struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	tracer trace.Tracer
}

func (p *Publisher) PublishEvent(ctx context.Context, id, subject string, payload []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, span := p.tracer.Start(ctx, "publish-nats-event")
		defer span.End()
		msg := nats.Msg{
			Subject: subject,
			Data:    payload,
			Header: nats.Header{
				"Nats-Msg-Id": []string{id},
			},
		}
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))
		ack, err := p.js.PublishMsg(&msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("failed to publish event: %w", err)
		}

		if ack == nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("failed to publish event: nil puback")
		}
		span.SetStatus(codes.Ok, "event published")
		return nil
	}
}

func (p *Publisher) EnsureStream(subjects []string) error {
	_, err := p.js.StreamInfo(StreamName)
	if err == nil {
		return nil
	}

	_, err = p.js.AddStream(&nats.StreamConfig{
		Name:       StreamName,
		Subjects:   subjects,
		Storage:    nats.FileStorage,
		Retention:  nats.WorkQueuePolicy,
		MaxAge:     7 * 24 * time.Hour,
		Replicas:   1,
		Duplicates: 10 * time.Minute,
	})
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	return nil
}

func (p *Publisher) Close() {
	if p.nc != nil {
		p.nc.Close()
	}
}

func NewPublisher(natsURL string, tracer trace.Tracer) (*Publisher, error) {
	nc, err := nats.Connect(
		natsURL,
		nats.Name("water-system"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	return &Publisher{
		nc:     nc,
		js:     js,
		tracer: tracer,
	}, nil
}
