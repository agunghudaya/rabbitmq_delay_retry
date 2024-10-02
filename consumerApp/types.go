package main

// Payload represents the incoming message structure from RabbitMQ
type Payload struct {
	Data Data `json:"data"`
}

// Data represents the actual message content
type Data struct {
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

// ExternalServiceResponse represents the response from the external service
type ExternalServiceResponse struct {
	Success bool `json:"success"`
}
