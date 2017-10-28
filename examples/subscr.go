package main

import (
	"flag"
	"log"

	"github.com/karamani/amqptool"
)

var amqpCnnArg string

func init() {
	flag.StringVar(&amqpCnnArg, "cnn", "", "amqp connection string")
	flag.Parse()
}

func main() {
	if err := amqptool.NewAMQPSubscriber(amqpCnnArg, "amqptooltst", 0).Process(handleOneMessage); err != nil {
		log.Fatalln(err.Error())
	}
}

func handleOneMessage(msg []byte) error {
	log.Println(string(msg))
	return nil
}
