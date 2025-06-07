package interfaces

import "net/http"

// IReceiver defines the interface for receiving and processing files
type IReceiver interface {
	HandleUpload(w http.ResponseWriter, r *http.Request)
}
