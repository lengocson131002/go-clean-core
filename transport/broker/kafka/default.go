package kafka

import "github.com/lengocson131002/go-clean/pkg/logger/logrus"

var (
	DefaultKafkaBroker = "127.0.0.1:9092"
	DefaultLogger      = logrus.NewLogrusLogger()
)
