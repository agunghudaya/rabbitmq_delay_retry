package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Call the external service
func callExternalService() (*ExternalServiceResponse, error) {

	// Make an HTTP POST request to the external service
	resp, err := http.Get("http://externalService:8081/api/transaction")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body into the expected structure
	var serviceResponse ExternalServiceResponse
	if err := json.Unmarshal(body, &serviceResponse); err != nil {
		return nil, err
	}

	return &serviceResponse, nil
}
