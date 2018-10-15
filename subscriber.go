package amqptool

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Subscriber contains data for subscribing to messages from the queue.
type Subscriber struct {
	Connection    string
	Queue         string
	Exchange      string
	PrefetchCount int
	QueueOpt      QueueOpt
	ConsumeOpt    ConsumeOpt
}

// NewSubscriber create & init the Subscriber struct.
func NewSubscriber(cnn, queue string) *Subscriber {
	return &Subscriber{
		Connection: cnn,
		Queue:      queue,
		QueueOpt: QueueOpt{
			Durable: true,
		},
	}
}

// SetExchange set the Exchange field.
func (s *Subscriber) SetExchange(ex string) *Subscriber {
	s.Exchange = ex
	return s
}

// SetPrefetchCount set the PrefetchCount field.
func (s *Subscriber) SetPrefetchCount(n int) *Subscriber {
	s.PrefetchCount = n
	return s
}

// AddQueueArg add argument to QueueOpt field.
func (s *Subscriber) AddQueueArg(key string, value interface{}) *Subscriber {
	if s.QueueOpt.Args == nil {
		s.QueueOpt.Args = make(map[string]interface{})
	}
	s.QueueOpt.Args[key] = value
	return s
}

// Process starts a message loop.
func (s *Subscriber) Process(h func([]byte) error) error {

	conn, err := amqp.Dial(s.Connection)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ %s", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel %s", err.Error())
	}
	defer ch.Close()

	if err := ch.Qos(s.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("qos error %s", err.Error())
	}

	if _, err := ch.QueueDeclare(
		s.Queue,
		s.QueueOpt.Durable,
		s.QueueOpt.AutoDelete,
		s.QueueOpt.Exclusive,
		s.QueueOpt.NoWait,
		amqp.Table(s.QueueOpt.Args)); err != nil {

		return fmt.Errorf("failed to declare a queue %s", err.Error())
	}

	if len(s.Exchange) > 0 {
		if err := ch.QueueBind(s.Queue, "", s.Exchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind a queue %s", err.Error())
		}
	}

	msgs, err := ch.Consume(
		s.Queue,
		s.ConsumeOpt.Consumer,
		s.ConsumeOpt.AutoAck,
		s.ConsumeOpt.Exclusive,
		s.ConsumeOpt.NoLocal,
		s.ConsumeOpt.NoWait,
		amqp.Table(s.ConsumeOpt.Args),
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer %s", err.Error())
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			go func(d amqp.Delivery) {
				defer func() {
					d.Ack(false)
				}()
				h(d.Body)
			}(d)
		}

		forever <- true
	}()

	select {
	case err = <-ch.NotifyClose(make(chan *amqp.Error)):
		err = fmt.Errorf("error from amqp.NotifyClose %s", err.Error())
	case <-forever:
		err = nil
	}

	return err
}
