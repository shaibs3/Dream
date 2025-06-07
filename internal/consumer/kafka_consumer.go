package kafkaConsumer

import (
	"encoding/json"
	"fmt"
	"log"

	"dream/internal/models"
	"dream/internal/parser"
	"dream/internal/types"

	"github.com/Shopify/sarama"
)

// KafkaConsumer implements the IConsumer interface for Kafka
type KafkaConsumer struct {
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	broker            string
	topic             string
	stopChan          chan struct{}
	storage           ProcessStorage
}

// NewKafkaConsumer creates a new KafkaConsumer
func NewKafkaConsumer(broker string, topic string, storage ProcessStorage) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	log.Printf("KafkaConsumer: Creating consumer for topic %s on broker %s", topic, broker)
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		broker:   broker,
		topic:    topic,
		stopChan: make(chan struct{}),
		storage:  storage,
	}, nil
}

// Start begins consuming messages from Kafka
func (kc *KafkaConsumer) Start() error {
	log.Printf("KafkaConsumer: Starting consumer for topic %s on broker %s", kc.topic, kc.broker)
	partitionConsumer, err := kc.consumer.ConsumePartition(kc.topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("KafkaConsumer: Failed to create partition consumer: %v", err)
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}
	kc.partitionConsumer = partitionConsumer

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("KafkaConsumer: Received message: %s", string(msg.Value))

				// Unmarshal the message
				var req types.MessageRequest
				if err := json.Unmarshal(msg.Value, &req); err != nil {
					log.Printf("KafkaConsumer: Error unmarshaling message: %v", err)
					continue
				}

				// Select parser using factory
				p, err := parser.GetParser(req)
				if err != nil {
					log.Printf("KafkaConsumer: Parser error: %v", err)
					continue
				}

				// Parse the command output
				entries, err := p.Parse(req.CommandOutput)
				if err != nil {
					log.Printf("KafkaConsumer: Error parsing command output: %v", err)
					continue
				}

				log.Printf("KafkaConsumer: Parsed %d process entries", len(entries))

				// Prepare user model
				user := &models.User{
					MachineID:   req.MachineID,
					MachineName: req.MachineName,
					OSVersion:   req.OSVersion,
					UserName:    req.UserName,
					UserID:      req.UserID,
				}

				// Convert entries to model structs
				var osType string
				var processModels []interface{}
				switch req.OSVersion {
				case "Windows 10":
					osType = "windows"
					for _, e := range entries {
						processModels = append(processModels, e.(parser.WindowsProcess))
					}
				default:
					osType = "linux"
					for _, e := range entries {
						processModels = append(processModels, e.(parser.LinuxProcess))
					}
				}

				log.Printf("KafkaConsumer: Storing user_id=%s, machine_id=%s, osType=%s, process_count=%d", user.UserID, user.MachineID, osType, len(processModels))
				if err := kc.storage.StoreUserAndProcesses(user, processModels, osType, req.Timestamp); err != nil {
					log.Printf("KafkaConsumer: DB store error: %v", err)
				}

			case err := <-partitionConsumer.Errors():
				log.Printf("KafkaConsumer: Error consuming message: %v", err)
			case <-kc.stopChan:
				log.Printf("KafkaConsumer: Stopping consumer")
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
