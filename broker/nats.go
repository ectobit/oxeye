package broker

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	defaultAckWait            = 60 * time.Second
	defaultReceiveChannelSize = 128
)

var _ Broker = (*NatsJetStream)(nil)

// NatsJetStream implements Broker interface for NATS JetStream broker.
// Exported field Debug can be used for debugging.
type NatsJetStream struct {
	c      nats.JetStreamContext
	config *NatsJetStreamConfig
	wg     sync.WaitGroup
	done   chan struct{}
	Debug  func(s string)
}

// NatsJetStreamConfig contains NatsJetStream configuration parameters.
type NatsJetStreamConfig struct {
	// Consume this subject
	ConsumeSubject string
	// Optional. If provided, queue group will be used.
	ConsumerGroup string
	// Produce into this subject
	ProduceSubject string
	// ReceiveChannelSize will prevent dropping messages caused by th slow consumer.
	ReceiveChannelSize int
	// Limit for pending messages. Default is -1 which means unlimited.
	PendingLimit int
	// How long to wait for ACK. If crossed, message will be redelivered. Default 60.s
	AckWait time.Duration
	// MaxRedeliveries defines how many times message will be redelivered if not acknowledged. Default 2.
	MaxRedeliveries uint8
}

// NewNatsJetStream creates new NATS JetStream broker implementing broker.Broker interface.
func NewNatsJetStream(client nats.JetStreamContext, config *NatsJetStreamConfig) *NatsJetStream {
	if config.AckWait == 0 {
		config.AckWait = defaultAckWait
	}

	if config.PendingLimit == 0 {
		config.PendingLimit = -1
	}

	if config.MaxRedeliveries == 0 {
		config.MaxRedeliveries = 2
	}

	if config.ReceiveChannelSize == 0 {
		config.ReceiveChannelSize = defaultReceiveChannelSize
	}

	return &NatsJetStream{ //nolint:exhaustruct
		c:      client,
		config: config,
		done:   make(chan struct{}),
		Debug:  func(string) {},
	}
}

// Sub implements broker.Broker interface.
func (b *NatsJetStream) Sub() (<-chan Message, error) { //nolint:funlen,cyclop,gocognit
	messages := make(chan Message)
	natsCh := make(chan *nats.Msg, b.config.ReceiveChannelSize)

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

	if err = sub.SetPendingLimits(b.config.PendingLimit, -1); err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	b.wg.Add(1)

	go func() {
		for {
			select {
			case msg := <-natsCh:
				messages <- Message{
					Data: msg.Data,
					Ack: func() {
						if err := msg.Ack(); err != nil {
							b.Debug(fmt.Sprintf("ack: %s", err))
						}
					},
					InProgress: func() {
						if err := msg.InProgress(); err != nil {
							b.Debug(fmt.Sprintf("in progress: %s", err))
						}
					},
				}
			case <-b.done:
				defer b.wg.Done()

				b.Debug("stopping consumer")

				if err := sub.Unsubscribe(); err != nil {
					b.Debug(fmt.Sprintf("unsubscribe: %s", err))
				}

				if err := sub.Drain(); err != nil {
					b.Debug(fmt.Sprintf("drain: %s", err))
				}

				close(natsCh)

				for range natsCh {
					<-natsCh
				}

				close(messages)

				return
			}
		}
	}()

	return messages, nil
}

// Pub implements broker.Broker interface.
func (b *NatsJetStream) Pub(data []byte) error {
	pub, err := b.c.Publish(b.config.ProduceSubject, data)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	b.Debug(fmt.Sprintf("publish stream: %s sequence: %d", pub.Stream, pub.Sequence))

	return nil
}

// Exit implements broker.Broker interface.
func (b *NatsJetStream) Exit() {
	close(b.done)
	b.wg.Wait()
}
