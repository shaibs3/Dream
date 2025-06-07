package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"dream/interfaces"
	"dream/kafkaConsumer"
	"dream/kafkaProducer"
	"dream/models"
	"dream/uploader"

	"github.com/joho/godotenv"

	"gorm.io/gorm"
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

func initKafkaConsumer(db *gorm.DB) (interfaces.IConsumer, error) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	consumer, err := kafkaConsumer.NewKafkaConsumer(kafkaBroker, "file-uploads", db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start consumer: %v", err)
	}

	return consumer, nil
}

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Initialize PostgreSQL connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	db, err := models.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. create the Producer
	producer, err := initKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	// 2. create the Consumer
	consumer, err := initKafkaConsumer(db)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			log.Printf("Error stopping consumer: %v", err)
		}
	}()

	// 3. Create the Application Api
	fileReceiver := uploader.NewFileReceiver(producer)
	http.HandleFunc("/upload", fileReceiver.HandleUpload)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
