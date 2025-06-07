package uploader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"dream/interfaces"
)

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

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	// Validate file type
	ext := filepath.Ext(header.Filename)
	if !isValidFileType(ext) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Send to Kafka
	err = fr.producer.SendMessage(content)
	if err != nil {
		http.Error(w, "Error sending message to Kafka", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename); err != nil {
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
