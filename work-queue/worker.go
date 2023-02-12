package main

import (
	"bytes"
	"log"
	"time"

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
	//task_queue ini adalah nama queue nya
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}

	//set berapa pekerjaan untuk queue nya
	err = ch.Qos(1, 0, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to set qos"))
	}

	//listening channel
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to consume queue"))
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Receive a message: %s", d.Body)
			//simulasi banyak queue
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Println("Done")
			d.Ack(false)
		}
	}()

	log.Printf("Waiting for message....")

	<-forever
}
