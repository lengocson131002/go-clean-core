package kafka

import (
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

const (
	CorrelationIdHeader = "correlationId"
)

type Marshaler interface {
	Marshal(topic string, msg *broker.Message) (*sarama.ProducerMessage, error)
}

type Unmarshaler interface {
	Unmarshal(*sarama.ConsumerMessage) (*broker.Message, error)
}

type Codec interface {
	Marshaler
	Unmarshaler
}

type DefaultMarshaler struct{}

func (DefaultMarshaler) Marshal(topic string, msg *broker.Message) (*sarama.ProducerMessage, error) {
	if len(msg.Headers) == 0 {
		msg.Headers = make(map[string]string)
	}

	correlationId, ok := msg.Headers[CorrelationIdHeader]
	if !ok || len(correlationId) == 0 {
		correlationId = uuid.New().String()
		msg.Headers[CorrelationIdHeader] = correlationId
	}

	headers := []sarama.RecordHeader{}

	for key, value := range msg.Headers {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}

	return &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.ByteEncoder(msg.Body),
		Headers: headers,
	}, nil
}

func (DefaultMarshaler) Unmarshal(kafkaMsg *sarama.ConsumerMessage) (*broker.Message, error) {
	headers := make(map[string]string, len(kafkaMsg.Headers))

	for _, header := range kafkaMsg.Headers {
		headers[string(header.Key)] = string(header.Value)
	}

	return &broker.Message{
		Headers: headers,
		Body:    []byte(kafkaMsg.Value),
	}, nil
}
