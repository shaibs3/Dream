package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestUploadFlow(t *testing.T) {
	// Start Docker Compose
	cmd := exec.Command("docker-compose", "up", "-d")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start Docker Compose: %v", err)
	}
	defer func() {
		cmd := exec.Command("docker-compose", "down")
		if err := cmd.Run(); err != nil {
			t.Logf("Error stopping Docker Compose: %v", err)
		}
	}()

	// Wait for services to be ready
	time.Sleep(10 * time.Second)

	// Create a test file
	testContent := []byte("test content")
	testFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Prepare the multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	if _, err := fileWriter.Write(testContent); err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}

	// Add metadata
	metadata := map[string]string{
		"user_id":     "123",
		"description": "test file",
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Failed to marshal metadata: %v", err)
	}
	if err := writer.WriteField("metadata", string(metadataBytes)); err != nil {
		t.Fatalf("Failed to write metadata field: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Send the request
	resp, err := http.Post("http://localhost:8080/upload", writer.FormDataContentType(), body)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Error closing response body: %v", err)
		}
	}()

	// Check response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Wait for message to be processed
	time.Sleep(5 * time.Second)

	// Create Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Start from the beginning
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		t.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			t.Logf("Error closing consumer: %v", err)
		}
	}()

	// Subscribe to the topic
	partitionConsumer, err := consumer.ConsumePartition("user-processes", 0, sarama.OffsetOldest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			t.Logf("Error closing partition consumer: %v", err)
		}
	}()

	// Read message with timeout
	select {
	case msg := <-partitionConsumer.Messages():
		// Verify message content
		var receivedMetadata map[string]interface{}
		if err := json.Unmarshal(msg.Value, &receivedMetadata); err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		// Verify metadata fields
		assert.Equal(t, "123", receivedMetadata["user_id"])
		assert.Equal(t, "test file", receivedMetadata["description"])
		assert.NotEmpty(t, receivedMetadata["file_path"])
		assert.NotEmpty(t, receivedMetadata["upload_time"])

	case err := <-partitionConsumer.Errors():
		t.Fatalf("Error consuming message: %v", err)

	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for Kafka message")
	}
}
