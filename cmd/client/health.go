package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func checkHealth(serverUrl string) error {
	healthEndpoint, err := url.JoinPath(serverUrl, "health")
	if err != nil {
		return fmt.Errorf("Error creating endpoints: %v", err)
	}

	resp, err := http.Get(healthEndpoint)
	if err != nil {
		return fmt.Errorf("Error making GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Health check failed with status %d", resp.StatusCode)
	}  

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %v", err)
	}

	log.Printf("Health response: %s", body)
	return nil
}
