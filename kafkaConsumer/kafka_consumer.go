package kafkaConsumer

import (
	"encoding/json"
	"fmt"
	"log"

	"dream/interfaces"
	"dream/parser"
	"dream/types"

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

				// Unmarshal the message
				var req types.MessageRequest
				if err := json.Unmarshal(msg.Value, &req); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}

				// Select parser using factory
				p, err := parser.GetParser(req)
				if err != nil {
					log.Printf("Parser error: %v", err)
					continue
				}

				// Parse the command output
				entries, err := p.Parse(req.CommandOutput)
				if err != nil {
					log.Printf("Error parsing command output: %v", err)
					continue
				}

				log.Printf("Parsed %d process entries", len(entries))

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
		if err := kc.partitionConsumer.Close(); err != nil {
			return fmt.Errorf("error closing partition consumer: %v", err)
		}
	}
	if kc.consumer != nil {
		if err := kc.consumer.Close(); err != nil {
			return fmt.Errorf("error closing consumer: %v", err)
		}
	}
	return nil
}
