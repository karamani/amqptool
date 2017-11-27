package amqptool

import (
	"fmt"

	"github.com/streadway/amqp"
)

type message struct {
	RoutingKey string
	Body       []byte
}

// Sender contains data for sending messages.
type Sender struct {
	Cnn             string
	Exchange        string
	ExchangeDeclare bool
	ExchangeOpt     ExchangeOpt

	errors    []chan error
	inC       chan message
	connected bool
}

// NewSender create & init the Sender struct.
func NewSender(cnn, exchange string) *Sender {
	return &Sender{
		Cnn:             cnn,
		Exchange:        exchange,
		ExchangeDeclare: true,
		ExchangeOpt: ExchangeOpt{
			Durable: true,
			Kind:    "direct",
		},
	}
}

// Connect creates a connection to the amqp-server and starts the message sending loop.
func (s *Sender) Connect() error {

	conn, err := amqp.Dial(s.Cnn)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	if s.ExchangeDeclare {
		if err := ch.ExchangeDeclare(
			s.Exchange,
			s.ExchangeOpt.Kind,
			s.ExchangeOpt.Durable,
			s.ExchangeOpt.Autodelete,
			s.ExchangeOpt.Internal,
			s.ExchangeOpt.NoWait,
			amqp.Table(s.ExchangeOpt.Args),
		); err != nil {
			ch.Close()
			conn.Close()
			return err
		}
	}

	go s.start(conn, ch)

	return nil
}

// Send sends a message to the exchange with a certain key.
func (s *Sender) Send(body []byte, routingKey string) error {

	if !s.connected {
		return fmt.Errorf("not connected")
	}

	s.inC <- message{
		RoutingKey: routingKey,
		Body:       body,
	}

	return nil
}

// NotifyError registers a listener for when the server sends a exception.
func (s *Sender) NotifyError(c chan error) chan error {
	s.errors = append(s.errors, c)
	return c
}

func (s *Sender) start(conn *amqp.Connection, ch *amqp.Channel) {

	defer func() {
		ch.Close()
		conn.Close()
		s.connected = false
		close(s.inC)
	}()

	s.inC = make(chan message)

	forever := make(chan error)

	go func() {
		for msg := range s.inC {
			if err := ch.Publish(
				s.Exchange,
				msg.RoutingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg.Body,
				},
			); err != nil {
				forever <- err
				return
			}
		}
		forever <- nil
	}()

	s.connected = true

	select {
	case err := <-ch.NotifyClose(make(chan *amqp.Error)):
		s.sendError(fmt.Errorf("NotifyClose %s", err.Error()))
	case err := <-forever:
		if err != nil {
			s.sendError(err)
		}
	}
}

func (s *Sender) sendError(e error) {
	for _, c := range s.errors {
		c <- e
	}
}
