package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"dream/interfaces"
	"dream/kafkaProducer"

	"github.com/Shopify/sarama"
)

var producer interfaces.IProducer

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

func consumeMessages() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Start from the beginning

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	consumer, err := sarama.NewConsumer([]string{kafkaBroker}, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("file-uploads", 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error consuming message: %v\n", err)
		}
	}
}

func main() {
	// Start consuming messages
	go func() {
		if err := consumeMessages(); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	// Your existing HTTP server code
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	if err := initKafkaProducer(); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	log.Println("Server starting on :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
