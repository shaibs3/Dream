package kafkaConsumer

// IConsumer defines the interface for a message consumer
type IConsumer interface {
	Start() error
	Stop() error
}

// ProcessStorage defines the interface for storing process data
type ProcessStorage interface {
	StoreUserAndProcesses(user interface{}, processes []interface{}, osType string, timestamp interface{}) error
}
