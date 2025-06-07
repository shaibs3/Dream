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

type Receiver struct {
	producer  types.IProducer
	validator Validator
}

func NewReceiver(producer types.IProducer, validator Validator) *Receiver {
	return &Receiver{
		producer:  producer,
		validator: validator,
	}
}

func (fr *Receiver) HandleReceive(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request to %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	log.Printf("Request body: %s", string(body))
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	// Parse the JSON request
	var messageReq types.MessageRequest
	if err := json.Unmarshal(body, &messageReq); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Use the injected validator
	if err := fr.validator.Validate(&messageReq); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send to Kafka
	err = fr.producer.SendMessage(body)
	if err != nil {
		log.Printf("Error sending message to Kafka: %v", err)
		http.Error(w, "Error sending message to Kafka", http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Request received for machine %s (%s)", messageReq.MachineName, messageReq.MachineID),
	}

	log.Printf("Successfully processed request for user_id=%s, machine_id=%s", messageReq.UserID, messageReq.MachineID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
