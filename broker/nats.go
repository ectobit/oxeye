package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.ectobit.com/lax"
)

const (
	consumeMessageChannelSize = 64
	defaultAckWait            = 60 * time.Second
)

var _ Broker = (*NatsJetStream)(nil)

// NatsJetStream implements Broker interface for NATS JetStream broker.
type NatsJetStream struct {
	c      nats.JetStreamContext
	config *NatsJetStreamConfig
	log    lax.Logger
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
func NewNatsJetStream(client nats.JetStreamContext, config *NatsJetStreamConfig, log lax.Logger) *NatsJetStream {
	if config.AckWait == 0 {
		config.AckWait = defaultAckWait
	}

	if config.MaxRedeliveries == 0 {
		config.MaxRedeliveries = 2
	}

	return &NatsJetStream{
		c:      client,
		config: config,
		log:    log,
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
			nats.ManualAck(), nats.AckWait(b.config.AckWait), nats.MaxDeliver(int(b.config.MaxRedeliveries)))
	} else {
		sub, err = b.c.ChanSubscribe(b.config.ConsumeSubject, natsCh, nats.ManualAck(),
			nats.AckWait(b.config.AckWait))
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
							b.log.Warn("ack", lax.Error(err))
						}
					},
				}
			case <-ctx.Done():
				b.log.Info("stopping consumer")

				if err := sub.Unsubscribe(); err != nil {
					b.log.Warn("unsubscribe: %w", lax.Error(err))
				}

				if err := sub.Drain(); err != nil {
					b.log.Warn("drain: %w", lax.Error(err))
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

	b.log.Debug("publish", lax.String("stream", pub.Stream), lax.Uint64("sequence", pub.Sequence))

	return nil
}
