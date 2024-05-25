package chatgpt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Thread struct {
	ID        string                 `json:"id"`
	Object    string                 `json:"object"`
	CreatedAt int                    `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (c *ChatGPT) CreateThread() (*Thread, error) {
	url := "https://api.openai.com/v1/threads"

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating thread: %w", err)
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error creating thread: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	var thread Thread
	err = json.NewDecoder(response.Body).Decode(&thread)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &thread, nil
}

func (c *ChatGPT) GetThread(threadID string) (*Thread, error) {
	url := "https://api.openai.com/v1/threads/" + threadID

	// Create a client
	client := &http.Client{}

	// Create a new GET request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	var thread Thread
	err = json.NewDecoder(response.Body).Decode(&thread)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &thread, nil
}
