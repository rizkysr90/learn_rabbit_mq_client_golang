package main

import (
	"context"
	"log"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

func main() {
	amqpConn, err := rabbitmq.New()
	defer func() {
		if err := amqpConn.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	//
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		message := amqp091.Publishing{
			Headers: amqp091.Table{
				"sample": "value",
			},
			Body: []byte("Hello " + strconv.Itoa(i)),
		}
		err := amqpChannel.PublishWithContext(ctx, "notification", "email", false, false, message)
		if err != nil {
			log.Println(err)
		}
	}
}
