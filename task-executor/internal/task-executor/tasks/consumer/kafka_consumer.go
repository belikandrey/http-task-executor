package consumer

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sourcegraph/conc/pool"
	"http-task-executor/task-executor/internal/task-executor/config"
	"http-task-executor/task-executor/internal/task-executor/logger"
	"http-task-executor/task-executor/internal/task-executor/tasks"
	"strings"
)

const (
	withoutTimeout = -1
)

// KafkaConsumer represents kafka consumer that consumes messages from kafka and send to executor.
type KafkaConsumer struct {
	// executor - implementation of tasks.Executor.
	executor tasks.Executor
	// consumer - kafka consumer.
	consumer *kafka.Consumer
	// topic - kafka topic name.
	topic string
	// consumerGroup - kafka consumer group ID.
	consumerGroup string
	// logger - logger implementation.
	logger logger.Logger
	// pool - workers pool.
	pool *pool.Pool
}

// NewKafkaConsumer - creates new KafkaConsumer.
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

	newPool := pool.New().WithMaxGoroutines(config.ExecutorNumGoroutines)

	return &KafkaConsumer{executor: executor, consumer: consumer, topic: config.KafkaCfg.Topic, logger: logger, pool: newPool}, nil
}

// Start started receiving messages from kafka and send to executor.
func (k *KafkaConsumer) Start(ctx context.Context) {
	k.logger.Info("Starting consumer")
	for {
		select {
		case <-ctx.Done():
			k.logger.Info("Consumer context canceled")
			return
		default:
			k.consumeMessages()
		}
	}
}

func (k *KafkaConsumer) consumeMessages() {
	msg, err := k.consumer.ReadMessage(withoutTimeout)
	if err != nil {
		k.logger.Errorf("Error reading message: %v", err)
		return
	}
	if msg == nil {
		k.logger.Warn("Consumer ReadMessage message is nil")
		return
	}
	k.logger.Debugf("Consumer ReadMessage: %v", string(msg.Value))

	k.pool.Go(func() {
		k.executor.Execute(msg.Value)
	})
}

// Close shutting down kafka consumer properly.
func (k *KafkaConsumer) Close(cancel context.CancelFunc) error {
	k.logger.Info("Shutting down consumer")

	k.pool.Wait()
	cancel()
	if _, err := k.consumer.Commit(); err != nil {
		if err.(kafka.Error).Code() != kafka.ErrNoOffset {
			k.logger.Errorf("Error closing consumer: %v", err)
			return err
		}
	}
	return k.consumer.Close()
}
