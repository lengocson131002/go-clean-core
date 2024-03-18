package kafka

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
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
		errCount     int64
	)

	err := kBroker.Connect()
	if err != nil {
		b.Fail()
	}

	_, err = kBroker.Subscribe(requestTopic, func(ctx context.Context, e broker.Event) error {
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
		err = kBroker.Publish(context.Background(), replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			b.Error(err)
			b.Fail()
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	round := 500
	total := 10

	var totalTime int64
	var wg sync.WaitGroup
	wg.Add(total * round)
	for i := 0; i < round; i++ {
		go func() {
			for j := 0; j < total; j++ {
				req := KRequestType{
					Number: rand.Intn(100),
				}
				reqByte, err := json.Marshal(req)
				if err != nil {
					b.Error(err)
				}

				start := time.Now()
				b.Logf("Start time: %v", start)
				msg, err := kBroker.PublishAndReceive(context.Background(), requestTopic, &broker.Message{
					Headers: map[string]string{
						"correlationId": uuid.New().String(),
					},
					Body: reqByte,
				}, broker.WithPublishReplyToTopic(replyTopic),
					broker.WithReplyConsumerGroup("benchmark.test"))

				if err != nil {
					errCount++
					b.Errorf("benchmark error: %v", err)
				} else {
					var result map[string]float64
					json.Unmarshal(msg.Body, &result)
					expected := math.Pow(float64(req.Number), 2)
					if result["Result"] != expected {
						b.Errorf("Expected result: %v, got: %v", expected, result["Result"])
					}
					b.Logf("Message: %v. Expected result: %v, got: %v, Duration: %vms", string(msg.Headers[CorrelationIdHeader]), expected, result["Result"], time.Since(start).Milliseconds())
				}

				totalTime += time.Since(start).Milliseconds()
				wg.Done()
			}
		}()
	}

	wg.Wait()
	b.Logf("AVG: %vms", totalTime/int64(total*round))
	kBroker.Disconnect()

}

func TestPublishAndReceived(t *testing.T) {
	var (
		kBroker      = getKafkaBroker()
		requestTopic = "go.clean.test.benchmark.request"
		replyTopic   = "go.clean.test.benchmark.reply"
	)

	err := kBroker.Connect()
	if err != nil {
		t.Error(err)
	}

	_, err = kBroker.Subscribe(requestTopic, func(ctx context.Context, e broker.Event) error {
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
		err = kBroker.Publish(context.Background(), replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			t.Error(err)
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		t.Error(err)
	}

	req := KRequestType{
		Number: rand.Intn(100),
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	msg, err := kBroker.PublishAndReceive(context.Background(), requestTopic, &broker.Message{
		Body: reqByte,
	},
		broker.WithPublishReplyToTopic(replyTopic),
		broker.WithReplyConsumerGroup("benchmark.test"))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	t.Logf("%v", msg)

}
