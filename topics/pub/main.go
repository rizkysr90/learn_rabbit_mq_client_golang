package main

import (
	"context"
	"log"
	"os"
	"strings"

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
	var routingKey string

	if len(os.Args) > 1 {
		routingKey = os.Args[1]
	} else {
		routingKey = "anonymous.info"
	}
	var message string

	if len(os.Args) > 2 {
		message = strings.Join(os.Args[2:], " ")
	} else {
		message = "Hello World!"
	}
	channel.ExchangeDeclare("topic_logs", "topic", false, false,
		false, false, amqp091.NewConnectionProperties())

	channel.PublishWithContext(context.Background(), "topic_logs", routingKey,
		false, false, amqp091.Publishing{
			Body: []byte(message),
		},
	)
	log.Println("sent message : ", routingKey, message)
	conn.Close()
}
