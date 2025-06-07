package kafkaConsumer

import (
	"fmt"
	"log"

	"dream/interfaces"

	"github.com/Shopify/sarama"
)

// KafkaConsumer implements the IConsumer interface for Kafka
type KafkaConsumer struct {
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	broker            string
	topic             string
	stopChan          chan struct{}
}

// NewKafkaConsumer creates a new KafkaConsumer
func NewKafkaConsumer(broker string, topic string) (interfaces.IConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		broker:   broker,
		topic:    topic,
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins consuming messages from Kafka
func (kc *KafkaConsumer) Start() error {
	partitionConsumer, err := kc.consumer.ConsumePartition(kc.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}
	kc.partitionConsumer = partitionConsumer

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("Received message: %s\n", string(msg.Value))
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming message: %v\n", err)
			case <-kc.stopChan:
				return
			}
		}
	}()

	return nil
}

// Stop stops consuming messages
func (kc *KafkaConsumer) Stop() error {
	close(kc.stopChan)
	if kc.partitionConsumer != nil {
		kc.partitionConsumer.Close()
	}
	if kc.consumer != nil {
		kc.consumer.Close()
	}
	return nil
}
