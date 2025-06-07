package receiver

import "net/http"

// IReceiver defines the interface for receiving and processing files
type IReceiver interface {
	HandleReceive(w http.ResponseWriter, r *http.Request)
}
