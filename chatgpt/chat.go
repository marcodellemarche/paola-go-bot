package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CompletionRequest struct {
	Model    string              `json:"model"`
	Messages []CompletionMessage `json:"messages"`
}

type CompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var initialMessage = CompletionMessage{
	Role:    "system",
	Content: "You are a nice friend named Paola. You're from Italy. You have a strong personality. You always sign off with her name.",
}

func (c *ChatGPT) CreateCompletion(content string, context []CompletionMessage) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	message := CompletionMessage{Role: "user", Content: content}

	messages := []CompletionMessage{initialMessage}
	messages = append(messages, context...)
	messages = append(messages, message)

	// Create the completionReq
	completionReq := CompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	// Convert the completion to JSON
	completionJSON, err := json.Marshal(completionReq)
	if err != nil {
		return "", err
	}

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(completionJSON))
	if err != nil {
		return "", err
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return "", err
	}

	// Read the response
	var completionRes CompletionResponse
	err = json.NewDecoder(response.Body).Decode(&completionRes)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var answer string
	for _, choice := range completionRes.Choices {
		if choice.Message.Role == "assistant" {
			answer = choice.Message.Content
			break
		}
	}

	return answer, nil
}

type CompletionChoice struct {
	Index        int               `json:"index"`
	Message      CompletionMessage `json:"message"`
	Logprobs     interface{}       `json:"logprobs"`
	FinishReason string            `json:"finish_reason"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	ID                string             `json:"id"`
	Object            string             `json:"object"`
	Created           int                `json:"created"`
	Model             string             `json:"model"`
	SystemFingerprint string             `json:"system_fingerprint"`
	Choices           []CompletionChoice `json:"choices"`
	Usage             CompletionUsage    `json:"usage"`
}
