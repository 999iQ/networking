package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strings"
)

const (
	flushTimeout = 5000 // ms
)

type Producer struct {
	producer *kafka.Producer
}

// NewProducer - функция для создания нового продюсера
func NewProducer(addrs []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addrs, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating the Kafka producer: %s", err)
	}
	return &Producer{producer: p}, nil
}

// Produce - функция для отправки сообщений в канал
func (p *Producer) Produce(msg, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(msg),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("error sending (produce) message to Kafka: %s", err)
	}
	// обработка пришедших сообщений из канала
	e := <-kafkaChan
	switch ev := e.(type) {
	// отправлено успешно
	case *kafka.Message:
		return nil
	// не отправлено по какой-то причине
	case kafka.Error:
		return ev
	default:
		return fmt.Errorf("unexpected message type=%T", e)
	}
}

// Close - функция для закрытия продюсера
func (p *Producer) Close() {
	// Flush - блокирует закрытие продюсера, пока все сообщения не отправятся или не таймер
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
