package receiver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"dream/internal/types"
)

// FileReceiver implements the IReceiver interface
// Now depends on a Validator

type FileReceiver struct {
	producer  types.IProducer
	validator Validator
}

// NewFileReceiver creates a new FileReceiver with a Validator
func NewFileReceiver(producer types.IProducer, validator Validator) types.IReceiver {
	return &FileReceiver{
		producer:  producer,
		validator: validator,
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

	// Use the injected validator
	if err := fr.validator.Validate(&messageReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
