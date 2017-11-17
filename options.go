package amqptool

// QueueOpt contains options for the queue.
type QueueOpt struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

// ConsumeOpt contains options for the Consumer.
type ConsumeOpt struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]interface{}
}
