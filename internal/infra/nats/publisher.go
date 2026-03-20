package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const StreamName = "WATER_SYSTEM"

type Publisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func (p *Publisher) PublishEvent(ctx context.Context, id, subject string, payload []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		ack, err := p.js.PublishMsg(&nats.Msg{
			Subject: subject,
			Data:    payload,
			Header: nats.Header{
				"Nats-Msg-Id": []string{id},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}

		if ack == nil {
			return fmt.Errorf("failed to publish event: nil puback")
		}

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

func NewPublisher(natsURL string) (*Publisher, error) {
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
		nc: nc,
		js: js,
	}, nil
}
