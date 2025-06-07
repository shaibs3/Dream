package uploader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"dream/interfaces"
	"dream/types"
)

// UploadRequest represents the structure of the incoming request
type UploadRequest struct {
	MachineID   string    `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	OSVersion   string    `json:"os_version"`
	Timestamp   time.Time `json:"timestamp"`
	CommandType string    `json:"command_type"` // e.g. "ps", "tasklist"
	BlobURL     string    `json:"blob_url"`
}

// FileReceiver implements the IReceiver interface
type FileReceiver struct {
	producer interfaces.IProducer
}

// NewFileReceiver creates a new FileReceiver
func NewFileReceiver(producer interfaces.IProducer) interfaces.IReceiver {
	return &FileReceiver{
		producer: producer,
	}
}

// HandleUpload handles the file upload request
func (fr *FileReceiver) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	// Parse the JSON request
	var messageReq types.MessageRequest
	if err := json.Unmarshal(body, &messageReq); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if messageReq.MachineID == "" || messageReq.MachineName == "" || messageReq.CommandType == "" ||
		messageReq.UserName == "" || messageReq.UserID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Send to Kafka
	err = fr.producer.SendMessage(body)
	if err != nil {
		http.Error(w, "Error sending message to Kafka", http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Request received for machine %s (%s)", messageReq.MachineName, messageReq.MachineID),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// isValidFileType checks if the file type is allowed
func isValidFileType(ext string) bool {
	allowedTypes := map[string]bool{
		".txt":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
	}
	return allowedTypes[ext]
}
