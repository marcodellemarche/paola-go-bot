package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Tool struct {
	Type string `json:"type"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type TruncationStrategy struct {
	Type         string      `json:"type"`
	LastMessages interface{} `json:"last_messages"`
}

type Function struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Function Function `json:"function"`
	Type     string   `json:"type"`
}

type SubmitToolOutputs struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

type RequiredAction struct {
	SubmitToolOutputs SubmitToolOutputs `json:"submit_tool_outputs"`
}

type Run struct {
	ID                  string                 `json:"id"`
	Object              string                 `json:"object"`
	CreatedAt           int                    `json:"created_at"`
	AssistantID         string                 `json:"assistant_id"`
	ThreadID            string                 `json:"thread_id"`
	Status              string                 `json:"status"`
	StartedAt           int                    `json:"started_at"`
	ExpiresAt           interface{}            `json:"expires_at"`
	CancelledAt         interface{}            `json:"cancelled_at"`
	FailedAt            interface{}            `json:"failed_at"`
	CompletedAt         int                    `json:"completed_at"`
	LastError           interface{}            `json:"last_error"`
	Model               string                 `json:"model"`
	Instructions        interface{}            `json:"instructions"`
	Tools               []Tool                 `json:"tools"`
	Metadata            map[string]interface{} `json:"metadata"`
	IncompleteDetails   interface{}            `json:"incomplete_details"`
	Usage               Usage                  `json:"usage"`
	Temperature         float64                `json:"temperature"`
	TopP                float64                `json:"top_p"`
	MaxPromptTokens     int                    `json:"max_prompt_tokens"`
	MaxCompletionTokens int                    `json:"max_completion_tokens"`
	TruncationStrategy  TruncationStrategy     `json:"truncation_strategy"`
	ResponseFormat      string                 `json:"response_format"`
	ToolChoice          string                 `json:"tool_choice"`
	RequiredAction      RequiredAction         `json:"required_action"`
}

type Assistant struct {
	AssistantID string `json:"assistant_id"`
}

func (c *ChatGPT) CreateRun(threadID string) (*Run, error) {
	url := "https://api.openai.com/v1/threads/" + threadID + "/runs"

	// Create the assistant
	assistant := Assistant{
		AssistantID: c.assistantID,
	}

	// Convert the assistant to JSON
	assistantJSON, err := json.Marshal(assistant)
	if err != nil {
		return nil, fmt.Errorf("error creating run: %w", err)
	}

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(assistantJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating run: %w", err)
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error creating run: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	var run Run
	err = json.NewDecoder(response.Body).Decode(&run)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &run, nil
}

type threadInput struct {
	Messages []messageInput `json:"messages"`
}

type threadAndRunInput struct {
	AssistantID string      `json:"assistant_id"`
	Thread      threadInput `json:"thread"`
}

func (c *ChatGPT) CreateThreadAndRun(content string) (*Run, error) {
	url := "https://api.openai.com/v1/threads/runs"

	// Create the input thread and run
	threadAndRunInput := threadAndRunInput{
		AssistantID: c.assistantID,
		Thread: threadInput{
			Messages: []messageInput{
				{
					Role:    "user",
					Content: content,
				},
			},
		},
	}

	// Convert the runThread to JSON
	runThreadJSON, err := json.Marshal(threadAndRunInput)
	if err != nil {
		return nil, fmt.Errorf("error creating thread and run: %w", err)
	}

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(runThreadJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating thread and run: %w", err)
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error creating thread and run: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	var run Run
	err = json.NewDecoder(response.Body).Decode(&run)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &run, nil
}

func (c *ChatGPT) GetRun(threadID string, runID string) (*Run, error) {
	url := "https://api.openai.com/v1/threads/" + threadID + "/runs/" + runID

	// Create a client
	client := &http.Client{}

	// Create a new GET request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the headers
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
	var run Run
	err = json.NewDecoder(response.Body).Decode(&run)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return &run, nil
}

type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

type ToolOutputs struct {
	ToolOutputs []ToolOutput `json:"tool_outputs"`
}

func (c *ChatGPT) SubmitToolOutputs(threadID string, runID string, toolCallID string, output string) error {
	url := "https://api.openai.com/v1/threads/" + threadID + "/runs/" + runID + "/submit_tool_outputs"

	// Create the tool outputs
	toolOutputs := ToolOutputs{
		ToolOutputs: []ToolOutput{
			{
				ToolCallID: toolCallID,
				Output:     output,
			},
		},
	}

	// Convert the tool outputs to JSON
	toolOutputsJSON, err := json.Marshal(toolOutputs)
	if err != nil {
		return fmt.Errorf("error creating tool outputs: %w", err)
	}

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(toolOutputsJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}

func (c *ChatGPT) CancelRun(threadID string, runID string) error {
	url := "https://api.openai.com/v1/threads/" + threadID + "/runs/" + runID + "/cancel"

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}

	// Check the response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	fmt.Println("Request successful")
	return nil
}
