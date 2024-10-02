package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Function to connect to RabbitMQ
func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// Function to declare the queue
func declareQueue(ch *amqp.Channel) (*amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"my_queue", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}
	return &queue, nil // Return a pointer to the queue
}

// Function to consume messages
func consumeMessages(ch *amqp.Channel) {
	q, err := declareQueue(ch)
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
		log.Printf("NEW message!! Body: %s", msg.Body)

		var message Payload
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}

		// Call the external service and check the response
		serviceResponse, err := callExternalService()
		if err != nil {
			log.Printf("Error calling external service: %s", err)
			continue
		}

		if !serviceResponse.Success {

			// Calculate the new delay using the Fibonacci sequence
			newDelay := nextFibonacci(message.Data.Delay)
			log.Printf("External service call failed for message: %s, with new delay: %d. republishing...", message.Data.Message, newDelay)

			// Republish the message with the new delay
			err = publishMessage(ch, message.Data.Message, newDelay)
			if err != nil {
				log.Printf("Error republishing message: %s", err)
			}
		} else {
			log.Printf("Received message: %s. Processed successfully.", message.Data.Message)
		}
	}
}

// Function to publish messages to RabbitMQ
func publishMessage(ch *amqp.Channel, message string, delay int) error {
	// Prepare the message body
	bodyData := map[string]interface{}{
		"data": map[string]interface{}{
			"message": message,
			"delay":   delay,
		},
	}

	body, err := json.Marshal(bodyData) // Marshal the bodyData to JSON
	if err != nil {
		return err
	}

	// Prepare the message properties and delay
	err = ch.Publish(
		"rabbitmq_consistent_hash_exchange", // Exchange
		"1",                                 // Routing key
		false,                               // Mandatory
		false,                               // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Headers: amqp.Table{
				"x-delay": int32(delay * 1000), // Set delay in milliseconds
			},
			Body: body,
		})

	return err
}

func nextFibonacci(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}

	a, b := 0, 1
	for b <= n {
		a, b = b, a+b
	}
	return b // Return the next Fibonacci number
}
