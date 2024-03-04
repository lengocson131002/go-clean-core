package broker

// Package broker is an interface used for asynchronous messaging

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...BrokerOption) error
	Options() BrokerOptions
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	PublishAndReceive(topic string, m *Message, opts ...PublishOption) (*Message, error)
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error

// Message is a message send/received from the broker.
type Message struct {
	Headers map[string]string
	Body    []byte
}

// Event is given to a subscription handler for processing.
type Event interface {
	// return event's topic
	Topic() string

	// return event's message
	Message() *Message

	// mark event as processed
	Ack() error

	// return error if event has error occurred
	Error() error
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}
