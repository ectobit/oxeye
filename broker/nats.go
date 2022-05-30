package broker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	consumeMessageChannelSize = 64
	defaultAckWait            = 60 * time.Second
)

var _ Broker = (*NatsJetStream)(nil)

// NatsJetStream implements Broker interface for NATS JetStream broker.
// Exported field Debug can be used for debugging.
type NatsJetStream struct {
	c      nats.JetStreamContext
	config *NatsJetStreamConfig
	Debug  io.StringWriter
}

// NatsJetStreamConfig contains NatsJetStream configuration parameters.
type NatsJetStreamConfig struct {
	// Consume this subject
	ConsumeSubject string
	// Optional. If provided, queue group will be used.
	ConsumerGroup string
	// Produce into this subject
	ProduceSubject string
	// How long to wait for ACK. If crossed, message will be redelivered. Default 60.s
	AckWait time.Duration
	// MaxRedeliveries defines how many times message will be redelivered if not acknowledged. Default 2.
	MaxRedeliveries uint8
}

// NewNatsJetStream creates new NATS JetStream broker.
func NewNatsJetStream(client nats.JetStreamContext, config *NatsJetStreamConfig) *NatsJetStream {
	if config.AckWait == 0 {
		config.AckWait = defaultAckWait
	}

	if config.MaxRedeliveries == 0 {
		config.MaxRedeliveries = 2
	}

	return &NatsJetStream{
		c:      client,
		config: config,
		Debug:  io.Discard.(io.StringWriter), //nolint:errcheck
	}
}

// Sub subscribes to broker and returns a channel to receive messages.
func (b *NatsJetStream) Sub(ctx context.Context) (<-chan Message, error) {
	messages := make(chan Message)
	natsCh := make(chan *nats.Msg, consumeMessageChannelSize)

	var sub *nats.Subscription

	var err error

	if b.config.ConsumerGroup != "" {
		sub, err = b.c.ChanQueueSubscribe(b.config.ConsumeSubject, b.config.ConsumerGroup, natsCh,
			nats.ManualAck(), nats.AckWait(b.config.AckWait), nats.MaxDeliver(int(b.config.MaxRedeliveries)),
			nats.DeliverNew())
	} else {
		sub, err = b.c.ChanSubscribe(b.config.ConsumeSubject, natsCh, nats.ManualAck(),
			nats.AckWait(b.config.AckWait), nats.DeliverNew())
	}

	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	go func() {
		for {
			select {
			case msg := <-natsCh:
				messages <- Message{
					Data: msg.Data,
					Ack: func() {
						if err := msg.Ack(); err != nil {
							b.Debug.WriteString(fmt.Sprintf("ack: %s", err))
						}
					},
					InProgress: func() {
						if err := msg.InProgress(); err != nil {
							b.Debug.WriteString(fmt.Sprintf("in progress: %s", err))
						}
					},
				}
			case <-ctx.Done():
				b.Debug.WriteString("stopping consumer")

				if err := sub.Unsubscribe(); err != nil {
					b.Debug.WriteString(fmt.Sprintf("unsubscribe: %s", err))
				}

				if err := sub.Drain(); err != nil {
					b.Debug.WriteString(fmt.Sprintf("drain: %s", err))
				}

				close(natsCh)
				close(messages)

				return
			}
		}
	}()

	return messages, nil
}

// Pub synchronously publishes a message to broker.
func (b *NatsJetStream) Pub(data []byte) error {
	pub, err := b.c.Publish(b.config.ProduceSubject, data)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	b.Debug.WriteString(fmt.Sprintf("publish stream: %s sequence: %d", pub.Stream, pub.Sequence))

	return nil
}
