package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

type KafkaBrokerConfig struct {
	Addresses []string

	SASLEnabled   bool
	SASLUser      string
	SASLPassword  string
	SASLAlgorithm string

	TLSEnabled        bool
	TLSSkipVerify     bool
	TLSClientCertFile string
	TLSClientKeyFile  string
	TLSCaCertFile     string
}

func createTLSConfiguration(certFile string, keyFile string, caFile string, skipVerify bool) (*tls.Config, error) {
	t := &tls.Config{
		InsecureSkipVerify: skipVerify,
	}

	if certFile != "" && keyFile != "" && caFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t.Certificates = []tls.Certificate{cert}
		t.RootCAs = caCertPool
	}

	return t, nil
}

func GetKafkaBroker(cfg *KafkaBrokerConfig, opts ...broker.BrokerOption) (broker.Broker, error) {
	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 1
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Metadata.Full = true

	// Config SASL
	if cfg.SASLEnabled {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = cfg.SASLUser
		conf.Net.SASL.Password = cfg.SASLPassword
		conf.Net.SASL.Handshake = true
		if cfg.SASLAlgorithm == "sha512" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		} else if cfg.SASLAlgorithm == "sha256" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256

		} else {
			return nil, fmt.Errorf("invalid SHA algorithm \"%s\": can be either \"sha256\" or \"sha512\"", cfg.SASLAlgorithm)
		}

	}

	// Config TLS
	if cfg.TLSEnabled {
		conf.Net.TLS.Enable = true

		tlsConfig, err := createTLSConfiguration(
			cfg.TLSClientCertFile,
			cfg.TLSClientKeyFile,
			cfg.TLSCaCertFile,
			cfg.TLSSkipVerify)

		if err != nil {
			return nil, fmt.Errorf("Failed to create TLS configuration: %w", err)
		}

		conf.Net.TLS.Config = tlsConfig

	}
	opts = append(
		opts,
		broker.WithBrokerAddresses(cfg.Addresses...),
		BrokerConfig(conf),
		ClusterConfig(conf))

	return NewKafkaBroker(opts...), nil

}
