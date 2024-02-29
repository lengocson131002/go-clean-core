package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

// consumerGroupHandler is the implementation of sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	logger  logger.Logger
	handler broker.Handler
	subopts broker.SubscribeOptions
	kopts   broker.BrokerOptions
	cg      sarama.ConsumerGroup
	sess    sarama.ConsumerGroupSession
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			ctx := context.Background()

			if !ok {
				h.logger.Warn(ctx, "[Kafka] message channel was closed")
				return nil
			}

			var m = broker.Message{}
			for _, header := range msg.Headers {
				m.Headers[string(header.Key)] = string(header.Value)
			}

			m.Body = []byte(msg.Value)
			p := &publication{m: &m, t: msg.Topic, km: msg, cg: h.cg, sess: session}

			err := h.handler(p)
			if err == nil && h.subopts.AutoAck {
				session.MarkMessage(msg, "")
			} else if err != nil {
				p.err = err
				errHandler := h.kopts.ErrorHandler
				if errHandler != nil {
					errHandler(p)
				} else {
					h.logger.Errorf(ctx, "[kafka] subscriber error: %v", err)
				}
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
