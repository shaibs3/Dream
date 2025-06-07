package kafkaProducer

import (
	"fmt"

	"dream/interfaces"

	"github.com/Shopify/sarama"
)

// KafkaProducer implements the IProducer interface for Kafka
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(broker string, topic string) (interfaces.IProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// SendMessage sends a message to the configured topic
func (kp *KafkaProducer) SendMessage(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kp.producer.SendMessage(msg)
	return err
}
