package main

import (
	"context"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

type FibonacciCall struct {
	channel       *amqp091.Channel
	callBack      *amqp091.Queue
	rpc           *amqp091.Queue
	correlationID string
}

func new(channel *amqp091.Channel, queue, rpc amqp091.Queue, correlationID string) *FibonacciCall {
	return &FibonacciCall{
		channel:       channel,
		callBack:      &queue,
		rpc:           &rpc,
		correlationID: correlationID,
	}

}
func (f *FibonacciCall) call(n int) {
	getReqBody := strconv.Itoa(n)
	if err := f.channel.PublishWithContext(context.Background(), "", f.rpc.Name, false, false, amqp091.Publishing{
		ReplyTo:       f.callBack.Name,
		DeliveryMode:  amqp091.Persistent,
		CorrelationId: f.correlationID,
		Body:          []byte(getReqBody),
	}); err != nil {
		panic(err)
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
	rpcQueue, err := channel.QueueDeclare("rpc_queue", true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	callbackQueue, err := channel.QueueDeclare("", true, false, true, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	callbackConsume, err := channel.ConsumeWithContext(context.Background(), callbackQueue.Name, callbackQueue.Name,
		true, false, false, false, amqp091.NewConnectionProperties())
	if err != nil {
		panic(err)
	}
	correlationID := uuid.NewString()
	fib := new(channel, callbackQueue, rpcQueue, correlationID)
	n := 5
	log.Printf("request fibonacci(%d)\n", n)
	var myChannel chan string
	defer close(myChannel)
	go func(channel chan string) {
		for message := range callbackConsume {
			if message.CorrelationId == correlationID {
				channel <- string(message.Body)

			}
		}
	}(myChannel)
	fib.call(5)
	log.Println("finished and waiting for the response")
	log.Println(correlationID)
	getRes := <-myChannel
	log.Println("Final", getRes)
}
