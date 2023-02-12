package main

import (
	"log"

	"github.com/andrisutanto/go-rabbitmq/broker"
	"github.com/pkg/errors"
)

func main() {
	conn, ch, err := broker.RabbitMQ()
	if err != nil {
		panic(err)
	}

	defer func() {
		ch.Close()
		conn.Close()
	}()

	//declare queue
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}

	//listening channel
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Receive a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for message....")

	<-forever
}
