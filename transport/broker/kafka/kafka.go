package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

var (
	CorrelationIdHeader = "correlationId"
	ReplyToTopicHeader  = "replyToTopic"
	RequestReplyTimeout = time.Second * 10
)

type kBroker struct {
	addrs []string

	c  sarama.Client        // broker connection client
	p  sarama.SyncProducer  // sync producer
	ap sarama.AsyncProducer // async producer

	cgs []sarama.ConsumerGroup

	connected bool
	scMutex   sync.Mutex
	opts      broker.BrokerOptions

	// request-reply patterns
	resps           map[string](chan *broker.Message) // response channels. key: message correlation ID, value: response body
	respSubscribers map[string]broker.Subscriber      // response subscribers. key: response topic, value: response subscriber
}

func NewKafkaBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.BrokerOptions{
		Context: context.Background(),
		Logger:  DefaultLogger, // using logrus logging by default
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}

	return &kBroker{
		addrs: cAddrs,
		opts:  options,
	}
}

type subscriber struct {
	k    *kBroker
	cg   sarama.ConsumerGroup
	t    string
	opts broker.SubscribeOptions
}

type publication struct {
	t    string
	err  error
	cg   sarama.ConsumerGroup
	km   *sarama.ConsumerMessage
	m    *broker.Message
	sess sarama.ConsumerGroupSession
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.sess.MarkMessage(p.km, "")
	return nil
}

func (p *publication) Error() error {
	return p.err
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	if err := s.cg.Close(); err != nil {
		return err
	}

	k := s.k
	k.scMutex.Lock()
	defer k.scMutex.Unlock()

	for i, cg := range k.cgs {
		if cg == s.cg {
			k.cgs = append(k.cgs[:i], k.cgs[i+1:]...)
			return nil
		}
	}

	return nil
}

func (k *kBroker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return DefaultKafkaBroker
}

func (k *kBroker) Connect() error {
	if k.connected {
		return nil
	}

	k.scMutex.Lock()
	if k.c != nil {
		k.scMutex.Unlock()
		return nil
	}
	k.scMutex.Unlock()

	pconfig := k.getBrokerConfig()
	// For implementation reasons, the SyncProducer requires
	// `Producer.Return.Errors` and `Producer.Return.Successes`
	// to be set to true in its configuration.
	pconfig.Producer.Return.Successes = true
	pconfig.Producer.Return.Errors = true

	c, err := sarama.NewClient(k.addrs, pconfig)
	if err != nil {
		return err
	}

	var (
		ap                   sarama.AsyncProducer
		p                    sarama.SyncProducer
		errChan, successChan = k.getAsyncProduceChan()
	)

	// Because error chan must require, so only error chan
	// If set the error chan, will use async produce
	// else use sync produce
	// only keep one client resource, is c variable
	if errChan != nil {
		ap, err = sarama.NewAsyncProducerFromClient(c)
		if err != nil {
			return err
		}
		// When the ap closed, the Errors() & Successes() channel will be closed
		// So the goroutine will auto exit
		go func() {
			for v := range ap.Errors() {
				errChan <- v
			}
		}()

		if successChan != nil {
			go func() {
				for v := range ap.Successes() {
					successChan <- v
				}
			}()
		}
	} else {
		p, err = sarama.NewSyncProducerFromClient(c)
		if err != nil {
			return err
		}
	}

	k.scMutex.Lock()
	k.c = c
	if p != nil {
		k.p = p
	}
	if ap != nil {
		k.ap = ap
	}
	k.cgs = make([]sarama.ConsumerGroup, 0)
	k.connected = true

	// request-reply pattern
	k.resps = make(map[string]chan *broker.Message)
	k.respSubscribers = make(map[string]broker.Subscriber)

	k.scMutex.Unlock()

	return nil
}

func (k *kBroker) Disconnect() error {
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	for _, consumer := range k.cgs {
		consumer.Close()
	}
	k.cgs = nil
	if k.p != nil {
		k.p.Close()
	}
	if k.ap != nil {
		k.ap.Close()
	}
	if err := k.c.Close(); err != nil {
		return err
	}
	k.connected = false

	// request-reply pattern
	k.resps = nil
	k.respSubscribers = nil

	return nil
}

func (k *kBroker) Init(opts ...broker.BrokerOption) error {
	for _, o := range opts {
		o(&k.opts)
	}
	var cAddrs []string
	for _, addr := range k.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}
	k.addrs = cAddrs
	return nil
}

func (k *kBroker) Options() broker.BrokerOptions {
	return k.opts
}

func (k *kBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	options := broker.PublishOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return k.sendMessage(topic, msg)
}

func (k *kBroker) PublishAndReceive(topic string, msg *broker.Message, opts ...broker.PublishOption) (*broker.Message, error) {
	options := broker.PublishOptions{
		ReplyToTopic: fmt.Sprintf("%s.reply", topic),
		Timeout:      RequestReplyTimeout,
	}

	for _, opt := range opts {
		opt(&options)
	}

	// Generate request Correlation ID if missed
	correlationId, ok := msg.Headers[CorrelationIdHeader]
	if !ok || len(correlationId) == 0 {
		if len(msg.Headers) == 0 {
			msg.Headers = make(map[string]string)
		}

		correlationId = uuid.New().String()
		msg.Headers[CorrelationIdHeader] = correlationId
	}

	var (
		replyTopic = options.ReplyToTopic
		timeout    = options.Timeout
	)

	// Subscribe for reply topic if didn't
	if _, ok := k.respSubscribers[replyTopic]; !ok {
		replySub, err := k.Subscribe(replyTopic, func(e broker.Event) error {
			if e.Message() == nil {
				return broker.EmptyRequestError{}
			}

			cId, correlationIdOk := e.Message().Headers[CorrelationIdHeader]
			if !correlationIdOk {
				return nil
			}

			msgChan, msgChanOk := k.resps[cId]
			if msgChanOk {
				msgChan <- e.Message()
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		k.respSubscribers[replyTopic] = replySub
	}

	k.resps[correlationId] = make(chan *broker.Message, 1)

	err := k.sendMessage(topic, msg)
	if err != nil {
		return nil, err
	}

	// wait for the response message
	msgChan := k.resps[correlationId]
	select {
	case body := <-msgChan:
		// remove processed channel
		delete(k.resps, correlationId)
		return body, nil
	case <-time.After(timeout):
		// remove processed channel
		delete(k.resps, correlationId)
		return nil, broker.RequestTimeoutResponse{
			Timeout: timeout,
		}
	}
}

func (k *kBroker) sendMessage(topic string, msg *broker.Message) error {
	kMsg := k.toKafkaMessage(topic, msg)
	if k.ap != nil {
		k.ap.Input() <- kMsg
		return nil
	} else if k.p != nil {
		_, _, err := k.p.SendMessage(kMsg)
		return err
	}
	return errors.New(`no connection resources available`)
}

func (k *kBroker) toKafkaMessage(topic string, msg *broker.Message) *sarama.ProducerMessage {
	kMsg := sarama.ProducerMessage{
		Topic:    topic,
		Value:    sarama.StringEncoder(msg.Body),
		Metadata: msg,
	}

	for k, v := range msg.Headers {
		kMsg.Headers = append(kMsg.Headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	return &kMsg
}

func (k *kBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		Group:   uuid.New().String(),
	}

	for _, o := range opts {
		o(&opt)
	}
	// we need to create a new client per consumer
	cg, err := k.getSaramaConsumerGroup(opt.Group)
	if err != nil {
		return nil, err
	}

	csHandler := &consumerGroupHandler{
		handler: handler,
		subopts: opt,
		kopts:   k.opts,
		cg:      cg,
		logger:  k.getLogger(),
	}

	ctx := context.Background()
	topics := []string{topic}
	go func() {
		for {
			select {
			case err := <-cg.Errors():
				if err != nil {
					k.getLogger().Errorf(ctx, "consumer error: %s", err)
				}
			default:
				err := cg.Consume(ctx, topics, csHandler)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					k.getLogger().Errorf(ctx, "consumer error: %s", err)
				}
			}
		}
	}()

	return &subscriber{
		k:    k,
		cg:   cg,
		opts: opt,
		t:    topic,
	}, nil
}

func (k *kBroker) getBrokerConfig() *sarama.Config {
	if c, ok := k.opts.Context.Value(brokerConfigKey{}).(*sarama.Config); ok {
		return c
	}
	return DefaultBrokerConfig
}

func (k *kBroker) getClusterConfig() *sarama.Config {
	if c, ok := k.opts.Context.Value(clusterConfigKey{}).(*sarama.Config); ok {
		return c
	}
	clusterConfig := DefaultClusterConfig

	// the oldest supported version is V0_10_2_0
	if !clusterConfig.Version.IsAtLeast(sarama.V0_10_2_0) {
		clusterConfig.Version = sarama.V0_10_2_0
	}

	clusterConfig.Consumer.Return.Errors = true
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	return clusterConfig
}

func (k *kBroker) getSaramaConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	config := k.getClusterConfig()
	cg, err := sarama.NewConsumerGroup(k.addrs, groupID, config)
	if err != nil {
		return nil, err
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.cgs = append(k.cgs, cg)
	return cg, nil
}

func (k *kBroker) getLogger() logger.Logger {
	logger := k.opts.Logger
	if logger == nil {
		logger = DefaultLogger
	}
	return logger
}

func (k *kBroker) getAsyncProduceChan() (chan<- *sarama.ProducerError, chan<- *sarama.ProducerMessage) {
	var (
		errors    chan<- *sarama.ProducerError
		successes chan<- *sarama.ProducerMessage
	)
	if c, ok := k.opts.Context.Value(asyncProduceErrorKey{}).(chan<- *sarama.ProducerError); ok {
		errors = c
	}
	if c, ok := k.opts.Context.Value(asyncProduceSuccessKey{}).(chan<- *sarama.ProducerMessage); ok {
		successes = c
	}
	return errors, successes
}

func (k *kBroker) String() string {
	return "kafka broker implementation"
}
