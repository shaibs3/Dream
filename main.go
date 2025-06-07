package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"dream/interfaces"
	"dream/kafkaConsumer"
	"dream/kafkaProducer"
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	err = producer.SendMessage(content)
	if err != nil {
		http.Error(w, "Error sending message to Kafka", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	if err := initKafkaProducer(); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	if err := initKafkaConsumer(); err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer consumer.Stop()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
