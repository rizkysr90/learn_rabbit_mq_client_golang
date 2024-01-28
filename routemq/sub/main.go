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
	res, err := channel.QueueDeclare("", false, false, true, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	queueName := res.Name
	severities := os.Args[1:]
	log.Println(severities)
	for _, severity := range severities {
		log.Println(string(severity))
		if err := channel.QueueBind(queueName, string(severity), "direct_logs",
			false, amqp091.NewConnectionProperties()); err != nil {
			panic(err)
		}
	}
	log.Println("waiting for the messages")
	consumer, err := channel.ConsumeWithContext(context.Background(), queueName, queueName,
		true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	for message := range consumer {
		log.Printf("%s : %s", message.RoutingKey, message.Body)
	}
}
