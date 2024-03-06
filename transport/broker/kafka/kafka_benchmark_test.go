package kafka

import (
	"encoding/json"
	"math"
	"math/rand"
	"testing"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/logger/logrus"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

func getKafkaBroker() broker.Broker {
	var config = &KafkaBrokerConfig{
		Addresses: []string{"localhost:9092"},
	}

	br, err := GetKafkaBroker(
		config,
	)

	if err != nil {
		panic(err)
	}

	return br
}

func getLogger() logger.Logger {
	return logrus.NewLogrusLogger()
}

type KRequestType struct {
	Number int
}

type KResponseType struct {
	Result float64
}

func BenchmarkKafka(b *testing.B) {
	var (
		kBroker      = getKafkaBroker()
		requestTopic = "go.clean.test.benchmark.request"
		replyTopic   = "go.clean.test.benchmark.reply"
	)

	err := kBroker.Connect()
	if err != nil {
		b.Fail()
	}

	_, err = kBroker.Subscribe(requestTopic, func(e broker.Event) error {
		msg := e.Message()
		if msg == nil {
			return broker.EmptyMessageError{}
		}

		var req KRequestType
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		result := KResponseType{
			Result: math.Pow(float64(req.Number), 2),
		}

		resultByte, err := json.Marshal(result)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		// pubish to response topic
		err = kBroker.Publish(replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			b.Error(err)
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := KRequestType{
			Number: rand.Intn(100),
		}
		reqByte, err := json.Marshal(req)
		if err != nil {
			b.Error(err)
		}

		err = kBroker.Publish(requestTopic, &broker.Message{
			Body: reqByte,
		},
			broker.WithPublishReplyToTopic(replyTopic),
			broker.WithReplyConsumerGroup("benchmark.test"))

		if err != nil {
			b.Logf("error: %v", err)
		}
	}
}
