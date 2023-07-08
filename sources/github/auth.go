package github

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// headers

// Accept: application/vnd.github.v3+json
// Authorization: Bearer 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

// GET https://api.github.com/rate_limit

func DO() {
	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new request
	req, err := http.NewRequest("GET", "https://api.github.com/rate_limit", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer TOKEN_VAR") // Replace TOKEN_VAR with your actual token

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	fmt.Println(string(body))
}
