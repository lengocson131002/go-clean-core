package broker

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/lengocson131002/go-clean/pkg/logger"
)

type BrokerOption func(*BrokerOptions)

type BrokerOptions struct {
	Context context.Context

	// underlying logger
	Logger logger.Logger

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler

	Addrs []string

	TLSConfig *tls.Config
}

func WithBrokerContext(ctx context.Context) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Context = ctx
	}
}

func WithBrokerAddresses(addrs ...string) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Addrs = addrs
	}
}

func WithLogger(log logger.Logger) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Logger = log
	}
}

func WithBrokerErrorHandler(handler Handler) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.ErrorHandler = handler
	}
}

func WithBrokerTLSConfig(t *tls.Config) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.TLSConfig = t
	}
}

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Context      context.Context
	Timeout      time.Duration
	ReplyToTopic string
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(opts *PublishOptions) {
		opts.Context = ctx
	}
}

func WithPublishTimeout(timeout time.Duration) PublishOption {
	return func(opts *PublishOptions) {
		opts.Timeout = timeout
	}
}

func WithPublishReplyToTopic(replyToTopic string) PublishOption {
	return func(opts *PublishOptions) {
		opts.ReplyToTopic = replyToTopic
	}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Context context.Context

	Group string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

func WithSubscribeContext(ctx context.Context) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Context = ctx
	}
}

func WithSubscribeGroup(gr string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Group = gr
	}
}

func WithSubscribeAutoAck(autoAck bool) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.AutoAck = autoAck
	}
}
