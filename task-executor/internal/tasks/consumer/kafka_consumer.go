package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sourcegraph/conc/pool"
	"strings"
	"task-executor/internal/config"
	"task-executor/internal/logger"
	"task-executor/internal/tasks"
)

const (
	withoutTimeout = -1
	numGoroutines  = 100
)

type KafkaConsumer struct {
	executor      tasks.Executor
	consumer      *kafka.Consumer
	topic         string
	consumerGroup string
	logger        logger.Logger
	stopped       bool
	pool          *pool.Pool
}

func NewKafkaConsumer(config *config.Config, executor tasks.Executor, logger logger.Logger) (*KafkaConsumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(config.KafkaCfg.Addresses, ","),
		"group.id":                 config.KafkaCfg.ConsumerGroup,
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       true,
		"enable.auto.offset.store": true,
		"auto.commit.interval.ms":  6000,
	}

	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	if err := consumer.Subscribe(config.KafkaCfg.Topic, nil); err != nil {
		return nil, err
	}

	newPool := pool.New().WithMaxGoroutines(numGoroutines)

	return &KafkaConsumer{executor: executor, consumer: consumer, topic: config.KafkaCfg.Topic, logger: logger, pool: newPool}, nil
}

func (k *KafkaConsumer) Start() {
	k.logger.Info("Starting consumer")
	for {
		if k.stopped {
			break
		}
		msg, err := k.consumer.ReadMessage(withoutTimeout)
		if err != nil {
			k.logger.Errorf("Error reading message: %v", err)
			continue
		}
		if msg == nil {
			k.logger.Warn("Consumer ReadMessage message is nil")
			continue
		}
		k.logger.Debugf("Consumer ReadMessage: %v", string(msg.Value))

		k.pool.Go(func() {
			k.executor.Execute(msg.Value)
		})
	}
}

func (k *KafkaConsumer) Close() error {
	k.pool.Wait()
	k.stopped = true
	if _, err := k.consumer.Commit(); err != nil {
		k.logger.Errorf("Error closing consumer: %v", err)
		return err
	}
	return k.consumer.Close()
}
