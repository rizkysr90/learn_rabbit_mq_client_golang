package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

func main() {
	var bodyMessage string
	// if len(os.Args) > 1 {
	// 	bodyMessage = strings.Join(os.Args[1:], " ")
	// }
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

	for i := 0; i < 5; i++ {
		bodyMessage = fmt.Sprintf("Hello_World_%v", i)
		message := amqp091.Publishing{
			Headers: amqp091.Table{
				"sample": "value",
			},
			Body:         []byte(bodyMessage + "i"),
			DeliveryMode: amqp091.Persistent,
		}

		err = amqpChannel.PublishWithContext(ctx, "notification", "email", false, false, message)
		if err != nil {
			log.Println(err)
		}
	}

}
