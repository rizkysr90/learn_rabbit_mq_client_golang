package main

import (
	"context"
	"log"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
func main() {
	conn, err := rabbitmq.New()
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		panic(err)
	}
	rpcQueue, err := channel.QueueDeclare("rpc_queue", true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	consume, err := channel.ConsumeWithContext(context.Background(), rpcQueue.Name, rpcQueue.Name, false, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	log.Println("awaiting rpc request")
	for message := range consume {
		n, err := strconv.Atoi(string(message.Body))
		if err != nil {
			panic(err)
		}
		log.Printf("fib %d \n", n)
		res := fib(n)
		resToStr := strconv.Itoa(res)
		channel.PublishWithContext(context.Background(), "", message.ReplyTo, false, false, amqp091.Publishing{
			CorrelationId: message.CorrelationId,
			Body:          []byte(resToStr),
		})
		channel.Ack(message.DeliveryTag, false)
	}

}
