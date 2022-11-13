// Package broker contains message broker abstraction and NATS JetStream broker implementation.
package broker

// Message contains data from the broker.
type Message struct {
	Data       []byte
	Ack        func()
	InProgress func()
}

// Broker defines common broker methods.
type Broker interface {
	// Sub subscribes to broker and returns a channel to receive messages.
	Sub() (<-chan Message, error)
	// Pub synchronously publishes a message to broker.
	Pub([]byte) error
	// Exit gracefully shuts down subscriber.
	Exit()
}
