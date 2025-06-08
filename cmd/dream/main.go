package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	kafkaConsumer "dream/internal/consumer"
	"dream/internal/models"
	kafkaProducer "dream/internal/producer"
	"dream/internal/receiver"

	"github.com/joho/godotenv"
)

func initKafkaProducer() (kafkaProducer.IProducer, error) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	producer, err := kafkaProducer.NewKafkaProducer(kafkaBroker, "user-processes")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}
	return producer, nil
}

func initKafkaConsumer(kafkaBroker string, storage kafkaConsumer.ProcessStorage) (kafkaConsumer.IConsumer, error) {
	consumer, err := kafkaConsumer.NewKafkaConsumer(
		kafkaBroker,
		os.Getenv("KAFKA_TOPIC"),
		storage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Kafka consumer: %v", err)
	}

	return consumer, nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := models.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize storage (postgress db)
	storage := models.NewPostgresStorage(db)

	// Initialize Kafka consumer
	consumer, err := initKafkaConsumer(os.Getenv("KAFKA_BROKER"), storage)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			log.Printf("Error stopping consumer: %v", err)
		}
	}()

	// Initialize Kafka producer
	producer, err := initKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	// Initialize HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize receiverHandler with validator and producer
	validator := &receiver.MessageRequestValidator{}
	receiverHandler := receiver.NewReceiver(producer, validator)

	// Set up routes
	http.HandleFunc("/upload", receiverHandler.HandleReceive)

	// Start the app server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
