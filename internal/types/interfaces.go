package types

import "net/http"

// IProducer defines the interface for a message producer
type IProducer interface {
	SendMessage(message []byte) error
}

// IReceiver defines the interface for a receiver/handler
type IReceiver interface {
	HandleUpload(w http.ResponseWriter, r *http.Request)
}
