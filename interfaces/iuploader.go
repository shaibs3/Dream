package interfaces

import "net/http"

// IUploader defines the interface for handling file uploads
type IUploader interface {
	HandleUpload(w http.ResponseWriter, r *http.Request)
}
