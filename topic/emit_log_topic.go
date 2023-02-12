package main

import (
	"fmt"
	"os"

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

	//declare exchange jika menggunakan pub sub / exchang
	err = ch.ExchangeDeclare("logs_topic", amqp091.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	//utk exchangenya (param pertama, dikasih nama exchange nya = "logs")
	err = ch.Publish("logs_topic", os.Args[1], false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(os.Args[2]),
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to publish message"))
	}

	fmt.Println("Send message:", os.Args[2])
}
