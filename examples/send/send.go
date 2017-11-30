package main

import (
	"flag"
	"log"

	"github.com/karamani/amqptool"
)

var cnnArg, keyArg, exArg string

func init() {
	flag.StringVar(&cnnArg, "cnn", "", "amqp connection string")
	flag.StringVar(&keyArg, "key", "", "routing key")
	flag.StringVar(&exArg, "ex", "", "exchange name")
	flag.Parse()
}

func main() {
	s := amqptool.NewSender(cnnArg, exArg)
	if err := s.Connect(); err != nil {
		log.Fatalln(err.Error())
	}

	go func() {
		select {
		case err := <-s.NotifyError(make(chan error)):
			log.Fatalln(err.Error())
		}
	}()

	if err := s.Send([]byte("testbody"), keyArg); err != nil {
		log.Fatalln(err.Error())
	}
}
