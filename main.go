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

var (
	producer interfaces.IProducer
	consumer interfaces.IConsumer
)

func initKafkaProducer() error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	var err error
	producer, err = kafkaProducer.NewKafkaProducer(kafkaBroker, "file-uploads")
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}
	return nil
}

func initKafkaConsumer() error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	var err error
	consumer, err = kafkaConsumer.NewKafkaConsumer(kafkaBroker, "file-uploads")
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}

	return nil
}

func main() {
	if err := initKafkaProducer(); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	if err := initKafkaConsumer(); err != nil {
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
