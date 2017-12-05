package amqptool

// Receiver contains data for receiving messages from the queue.
type Receiver struct{}

// NewReceiver create & init the Receiver struct.
func NewReceiver(cnn, queue string) *Receiver {
	return &Receiver{}
}
