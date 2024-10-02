package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Response struct to send JSON response
type Response struct {
	Success bool `json:"success"`
}

// oneInTwenty returns true with 1 in 20 probability
func oneInTwenty() bool {
	rand.Seed(time.Now().UnixNano()) // Seed the random generator
	num := rand.Intn(20) + 1         // Generate random number between 1 and 20
	return num == 1                  // Return true only if the number is 1
}

// transactionHandler handles the /api/transaction request
func transactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Create a response with the result from oneInTwenty
	response := Response{
		Success: oneInTwenty(),
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send the response as JSON
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Register the handler for the /api/transaction route
	http.HandleFunc("/api/transaction", transactionHandler)

	// Start the server on port 8080
	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
