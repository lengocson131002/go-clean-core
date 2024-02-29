package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

var (
	DefaultBrokerConfig  = sarama.NewConfig()
	DefaultClusterConfig = sarama.NewConfig()
)

// Config for publisher
type brokerConfigKey struct{}

// Config for subscriber
type clusterConfigKey struct{}

func BrokerConfig(c *sarama.Config) broker.BrokerOption {
	return setBrokerOption(brokerConfigKey{}, c)
}

func ClusterConfig(c *sarama.Config) broker.BrokerOption {
	return setBrokerOption(clusterConfigKey{}, c)
}

type subscribeContextKey struct{}

// SubscribeContext set the context for broker.SubscribeOption
func SubscribeContext(ctx context.Context) broker.SubscribeOption {
	return setSubscribeOption(subscribeContextKey{}, ctx)
}

type subscribeConfigKey struct{}

func SubscribeConfig(c *sarama.Config) broker.SubscribeOption {
	return setSubscribeOption(subscribeConfigKey{}, c)
}
