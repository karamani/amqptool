package amqptool

// ExchangeOpt contains options for the exchange.
type ExchangeOpt struct {
	Kind       string
	Durable    bool
	Autodelete bool
	Internal   bool
	NoWait     bool
	Args       map[string]interface{}
}

// QueueOpt contains options for the queue.
type QueueOpt struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

// ConsumeOpt contains options for the consumer.
type ConsumeOpt struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]interface{}
}
