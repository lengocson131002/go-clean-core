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
	ready   chan bool
	codec   Codec
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
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
				h.logger.Info(ctx, "[kafka consumer] message channel was closed")
				return nil
			}

			if msg == nil || len(msg.Value) == 0 {
				continue
			}

			m, err := h.codec.Unmarshal(msg)
			if err != nil {
				h.logger.Errorf(ctx, "[kafka consumer]: failed to unmarshal consumed message: %v", err)
				continue
			}

			p := &publication{m: m, t: msg.Topic, km: msg, cg: h.cg, sess: session}

			err = h.handler(ctx, p)
			if err == nil && h.subopts.AutoAck {
				session.MarkMessage(msg, "")
			} else if err != nil {
				p.err = err
				errHandler := h.kopts.ErrorHandler
				if errHandler != nil {
					errHandler(ctx, p)
				} else {
					h.logger.Errorf(ctx, "[kafka] subscriber error: %v", err)
				}
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
