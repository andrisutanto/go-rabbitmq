package main

import (
	"log"

	"github.com/andrisutanto/go-rabbitmq/broker"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
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
	//task_queue ini adalah nama queue nya
	//exclusive: jika true, maka saat tidak digunakan lagi akan dihapus
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}

	//declare exchange jika menggunakan pub sub / exchang
	err = ch.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	//saat exchange akan mengirimkan pesan, maka harus menggunakan binding
	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to bind queue"))
	}

	//listening channel
	//auto ack di set true
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to consume queue"))
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Receive a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for message....")

	<-forever
}
