package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"dream/kafka"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	producer *kafka.Producer
	maxSize  int64
}

func NewUploadHandler(producer *kafka.Producer, maxSize int64) *UploadHandler {
	log.Printf("Initializing upload handler with max size: %d bytes", maxSize)
	return &UploadHandler{
		producer: producer,
		maxSize:  maxSize,
	}
}

type UploadRequest struct {
	Metadata map[string]string `json:"metadata"`
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	// Get the file from the request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	log.Printf("Received file: %s, size: %d bytes", header.Filename, header.Size)

	// Get metadata from the request
	var req UploadRequest
	if err := json.Unmarshal([]byte(c.PostForm("metadata")), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata format"})
		return
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	// Check file size
	contentSize := int64(len(content))
	log.Printf("File content size: %d bytes, max allowed: %d bytes", contentSize, h.maxSize)
	if contentSize > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File too large: %d bytes (max: %d bytes)", contentSize, h.maxSize),
		})
		return
	}

	// Create message
	msg := &kafka.FileUploadMessage{
		FileName:    header.Filename,
		FileContent: content,
		Metadata:    req.Metadata,
		Timestamp:   time.Now(),
	}

	// Send to Kafka
	if err := h.producer.SendMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message to Kafka"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File uploaded successfully: %s", header.Filename),
	})
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
