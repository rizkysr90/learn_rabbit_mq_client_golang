package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func New() (*amqp091.Connection, error) {
	connection, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	log.Println("AMQP Dial is success")
	return connection, nil
}
