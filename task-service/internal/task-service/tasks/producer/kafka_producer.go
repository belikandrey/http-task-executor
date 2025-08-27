package producer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"http-task-executor/task-service/internal/task-service/config"
	"http-task-executor/task-service/internal/task-service/logger"
	"http-task-executor/task-service/internal/task-service/models"
	"http-task-executor/task-service/internal/task-service/tasks/mapper"
)

const flushTimeout = 5000

// TaskProducer represents message producer to Kafka.
type TaskProducer struct {
	// producer - kafka producer implementation
	producer *kafka.Producer
	// topic - name of kafka topic
	topic string
	// logger - implementation of common logger
	logger logger.Logger
}

// NewTaskProducer creates new instance of TaskProducer.
func NewTaskProducer(config *config.Config, logger logger.Logger) (*TaskProducer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.KafkaCfg.Addresses, ","),
	}

	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	return &TaskProducer{producer: producer, topic: config.KafkaCfg.Topic, logger: logger}, nil
}

// Produce produced send message to kafka.
func (p *TaskProducer) Produce(task *models.Task) error {
	message := mapper.MapTaskToKafkaTaskMessage(task)

	bytes, err := json.Marshal(message)
	if err != nil {
		p.logger.Errorf("Error marshalling message to bytes: %v", err)

		return err
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)

	err = p.producer.Produce(kafkaMsg, kafkaChan)
	if err != nil {
		p.logger.Errorf("Error producing message to Kafka: %v", err)

		return err
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		p.logger.Errorf("Error producing message to Kafka: %v", ev)

		return ev
	default:
		p.logger.Errorf("Unuexpected event type from kafka: %T", ev)

		return fmt.Errorf("unexpected event type=%T", ev)
	}
}

// Close closes producer.
func (p *TaskProducer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
