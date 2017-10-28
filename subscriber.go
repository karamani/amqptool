package amqptool

import (
	"fmt"

	"github.com/streadway/amqp"
)

// AMQPSubscriber contains data for subscribing to messages from the queue.
type AMQPSubscriber struct {
	CnnString     string
	Queue         string
	PrefetchCount int
}

func NewAMQPSubscriber(cnn, queue string, prefetchCount int) *AMQPSubscriber {
	return &AMQPSubscriber{CnnString: cnn, Queue: queue, PrefetchCount: prefetchCount}
}

func (pc *AMQPSubscriber) Process(h func([]byte) error) error {

	conn, err := amqp.Dial(pc.CnnString)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ %s", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel %s", err.Error())
	}
	defer ch.Close()

	err = ch.Qos(pc.PrefetchCount, 0, false)
	if err != nil {
		return fmt.Errorf("qos error %s", err.Error())
	}

	q, err := ch.QueueDeclarePassive(
		pc.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue %s", err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
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

				if err := h(d.Body); err != nil {
					// TODO
				}
			}(d)
		}

		forever <- true
	}()

	select {
	case err = <-ch.NotifyClose(make(chan *amqp.Error)):
		err = fmt.Errorf("NotifyClose %s", err.Error())
	case <-forever:
		err = nil
	}

	return err
}
