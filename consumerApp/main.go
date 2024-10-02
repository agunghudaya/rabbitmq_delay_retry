package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Data struct {
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

type Payload struct {
	Data Data `json:"data"`
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating a channel: %s", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"my_queue", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring a queue: %s", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Error registering a consumer: %s", err)
	}

	log.Printf("Waiting for messages. To exit press CTRL+C")

	// Infinite loop to process messages
	for msg := range msgs {
		log.Printf("Raw message body: %s", msg.Body)

		var message Payload
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}

		// Wait for the specified delay
		log.Printf("Received message: %s. After waiting for %d seconds...", message.Data.Message, message.Data.Delay)

	}
}
