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
	channel.ExchangeDeclare("direct_logs", "direct", false, false, false, false, amqp091.NewConnectionProperties())
	channel.PublishWithContext(context.Background(), "direct_logs", "severity", false, false, amqp091.Publishing{
		Body: []byte("testing directs log"),
	})
	log.Println("message already sent!")
	conn.Close()
}
