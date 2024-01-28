package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rizkysr90/learn_rabbit_mq_client_golang/pkg/rabbitmq"
)

func main() {
	amqpConn, err := rabbitmq.New()
	if err != nil {
		panic(err)
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	amqpChannel.Qos(1, 1, true)
	ctx := context.Background()
	emailConsumer, err := amqpChannel.ConsumeWithContext(ctx,
		"email", "consumer-email", false, false, false, false, nil)
	log.Println("Launch email consumer")
	if err != nil {
		panic(err)
	}
	for message := range emailConsumer {
		// dotCount := bytes.Count(message.Body, []byte{'.'})
		time.Sleep(5 * time.Second)
		println("Routing Key: " + message.RoutingKey)
		println("Body: " + string(message.Body))
		// splitMessage := strings.Split(string(message.Body), " ")
		// Specify the folder path
		folderPath := "./results"
		fileName := fmt.Sprintf("HelloWorld_%v.txt", time.Now().UTC().Unix())
		// Specify the file path
		filePath := filepath.Join(folderPath, fileName)
		// Create the folder if it doesn't exist
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
		content := fmt.Sprintf("CONTENT FROM : ./%v IN %v", filePath, time.Now().UTC().String())

		// Open the file in write-only mode, create it if it doesn't exist, truncate it if it does
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		amqpChannel.Ack(message.DeliveryTag, false)

		defer file.Close()

		// Write the string content to the file
		_, err = file.WriteString(content)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Println("File created successfully:", filePath)
	}

}
