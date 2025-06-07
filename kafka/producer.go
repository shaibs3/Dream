package kafka

import (
	"encoding/json"
	"time"

	"dream/config"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

type FileUploadMessage struct {
	FileName    string            `json:"file_name"`
	FileContent []byte            `json:"file_content"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.Producer.MaxRetries
	config.Producer.Retry.Backoff = time.Millisecond * 100
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    cfg.Topic,
	}, nil
}

func (p *Producer) SendMessage(msg *FileUploadMessage) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(jsonData),
	}

	_, _, err = p.producer.SendMessage(message)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
