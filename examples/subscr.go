package main

import (
	"flag"
	"log"

	"github.com/karamani/amqptool"
)

var cnnArg, queueArg, exArg string

func init() {
	flag.StringVar(&cnnArg, "cnn", "", "amqp connection string")
	flag.StringVar(&queueArg, "q", "", "queue name")
	flag.StringVar(&exArg, "ex", "", "exchange name")
	flag.Parse()
}

func main() {
	s := amqptool.NewSubscriber(cnnArg, queueArg).SetExchange(exArg)
	if err := s.Process(handleOneMessage); err != nil {
		log.Fatalln(err.Error())
	}
}

func handleOneMessage(msg []byte) error {
	log.Println(string(msg))
	return nil
}
