package main

import (
	"log"
)

func main() {
	// Connect to RabbitMQ
	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Start consuming messages
	consumeMessages(ch)
}
