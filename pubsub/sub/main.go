package main

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

func main() {
	conn, err := rabbitmq.New()
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare("", false, false, true, false, amqp091.Table{})
	if err != nil {
		panic(err)
	}
	channel.QueueBind(queue.Name, queue.Name, "logs", false, amqp091.NewConnectionProperties())
	log.Println("listening")
	delivery, err := channel.ConsumeWithContext(context.Background(), queue.Name, "",
		true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	for message := range delivery {
		log.Println(string(message.Body))
	}
}
