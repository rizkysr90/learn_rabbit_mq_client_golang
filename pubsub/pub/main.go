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
	channel.ExchangeDeclare("logs", "fanout", true, false, false, false, amqp091.NewConnectionProperties())
	channel.PublishWithContext(
		context.Background(), "logs", "", false, false, amqp091.Publishing{
			Body: []byte("HelloFromFanout"),
		},
	)
	log.Println("sent message")
}
