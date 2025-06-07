package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"dream/interfaces"
	"dream/kafkaConsumer"
	"dream/kafkaProducer"
	"dream/uploader"
)

func initKafkaProducer() (interfaces.IProducer, error) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	producer, err := kafkaProducer.NewKafkaProducer(kafkaBroker, "file-uploads")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}
	return producer, nil
}

func initKafkaConsumer() (interfaces.IConsumer, error) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	consumer, err := kafkaConsumer.NewKafkaConsumer(kafkaBroker, "file-uploads")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start consumer: %v", err)
	}

	return consumer, nil
}

func main() {
	producer, err := initKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	consumer, err := initKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			log.Printf("Error stopping consumer: %v", err)
		}
	}()

	// Create and use the file uploader
	fileUploader := uploader.NewFileUploader(producer)
	http.HandleFunc("/upload", fileUploader.HandleUpload)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
