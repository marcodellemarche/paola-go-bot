package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Text struct {
	Value       string        `json:"value"`
	Annotations []interface{} `json:"annotations"`
}

type Content struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Message struct {
	ID          string        `json:"id"`
	Object      string        `json:"object"`
	CreatedAt   int           `json:"created_at"`
	ThreadID    string        `json:"thread_id"`
	Role        string        `json:"role"`
	Content     []Content     `json:"content"`
	AssistantID string        `json:"assistant_id"`
	RunID       string        `json:"run_id"`
	Attachments []interface{} `json:"attachments"`
	Metadata    interface{}   `json:"metadata"`
}

type MessageList struct {
	Object string    `json:"object"`
	Data   []Message `json:"data"`
}

type messageInput struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *ChatGPT) CreateMessage(threadID string, input string) (*Message, error) {
	url := "https://api.openai.com/v1/threads/" + threadID + "/messages"

	// Create the message
	messageInput := messageInput{
		Role:    "user",
		Content: input,
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(messageInput)
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	var message Message
	err = json.NewDecoder(response.Body).Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &message, nil
}

func (c *ChatGPT) ListMessages(threadID string) (*MessageList, error) {
	url := "https://api.openai.com/v1/threads/" + threadID + "/messages"

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
	var messageList MessageList
	err = json.NewDecoder(response.Body).Decode(&messageList)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &messageList, nil
}
