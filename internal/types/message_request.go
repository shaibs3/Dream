package types

import "time"

// MessageRequest represents the structure of the incoming request
type MessageRequest struct {
	MachineID     string    `json:"machine_id"`
	MachineName   string    `json:"machine_name"`
	OSVersion     string    `json:"os_version"`
	Timestamp     time.Time `json:"timestamp"`
	CommandType   string    `json:"command_type"` // e.g. "ps", "tasklist"
	UserName      string    `json:"user_name"`
	UserID        string    `json:"user_id"`
	CommandOutput string    `json:"command_output"`
	Faculty       string    `json:"faculty"`
}
