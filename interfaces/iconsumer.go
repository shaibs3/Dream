package interfaces

// IConsumer defines the interface for a message consumer
type IConsumer interface {
	Start() error
	Stop() error
}
