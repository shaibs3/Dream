package producer

// IProducer defines the interface for a message producer
type IProducer interface {
	SendMessage(message []byte) error
}
