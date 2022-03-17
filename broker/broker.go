// Package broker contains message broker abstraction and NATS JetStream broker implementation.
package broker

import "context"

// Message contains data from the broker.
type Message struct {
	Data       []byte
	Ack        func()
	InProgress func()
}

// Broker defines common broker methods.
type Broker interface {
	Sub(context.Context) (<-chan Message, error)
	Pub([]byte) error
}
