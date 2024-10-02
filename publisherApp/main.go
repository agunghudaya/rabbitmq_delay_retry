package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

// MessagePayload struct to capture the JSON payload
type MessagePayload struct {
	Data struct {
		Message string `json:"message"`
		Delay   int    `json:"delay"`
	} `json:"data"`
}

// Connect to RabbitMQ and return the channel
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

// Function to bind the queue to the exchange
func bindQueue(ch *amqp.Channel) error {
	// Declare the exchange
	err := ch.ExchangeDeclare(
		"rabbitmq_consistent_hash_exchange", // Exchange name
		"x-delayed-message",                 // Exchange type
		true,                                // Durable
		false,                               // Auto-deleted when unused
		false,                               // Internal
		false,                               // No-wait
		amqp.Table{"x-delayed-type": []byte("direct")}, // Set the type of the delayed exchange
	)
	if err != nil {
		return err
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		"my_queue", // Queue name
		true,       // Durable
		false,      // Delete when unused
		false,      // Exclusive
		false,      // No-wait
		nil,        // Arguments
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange with an integer routing key
	err = ch.QueueBind(
		"my_queue",                          // The queue we want to bind
		"1",                                 // Example routing key (must be an integer as a string)
		"rabbitmq_consistent_hash_exchange", // Exchange name
		false,                               // No-wait
		nil,                                 // No arguments
	)
	return err
}

// Function to publish a message to RabbitMQ
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
			Body: body, // Use marshalled JSON as body
		})

	return err
}

// Echo handler for publishing the message
func handlePublishMessage(c echo.Context) error {

	payload := new(MessagePayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	message := payload.Data.Message
	delay := payload.Data.Delay

	log.Println("got message to publish with delay:", delay)

	// Connect to RabbitMQ
	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to RabbitMQ"})
	}
	defer conn.Close()
	defer ch.Close()

	// Bind the queue to the exchange
	err = bindQueue(ch)
	if err != nil {
		log.Fatalf("Failed to bind queue: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bind queue"})
	}

	// Publish the message with delay
	err = publishMessage(ch, message, delay)
	if err != nil {
		log.Fatalf("Failed to publish message: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to publish message"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "Message published",
		"message": message,
		"delay":   (time.Duration(delay) * time.Second).String(), // Convert the duration to string correctly
	})
}

func main() {
	// Echo instance
	e := echo.New()

	// Route to publish message
	e.POST("/api/publish", handlePublishMessage)

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
