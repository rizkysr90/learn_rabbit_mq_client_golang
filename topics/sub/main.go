package main

import (
	"context"
	"log"
	"os"

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
	channel.ExchangeDeclare("topic_logs", "topic", false, false,
		false, false, amqp091.NewConnectionProperties())
	res, err := channel.QueueDeclare("", false, false, true, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	binding_keys := os.Args[1:]
	if len(binding_keys) == 0 {
		log.Println("invalid binding key")
		return
	}
	for _, binding_key := range binding_keys {
		channel.QueueBind(res.Name, binding_key, "topic_logs", false, amqp091.NewConnectionProperties())
	}
	log.Println("waiting for logs")
	consumer, err := channel.ConsumeWithContext(context.Background(), res.Name, res.Name,
		true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	for message := range consumer {
		log.Printf("%s : %s", message.RoutingKey, message.Body)
	}
}
