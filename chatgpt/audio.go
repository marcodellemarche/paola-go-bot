package chatgpt

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (c *ChatGPT) Transcript(filename string) (string, error) {
	url := "https://api.openai.com/v1/audio/transcriptions"

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Add the other fields
	fieldWriter, err := multiPartWriter.CreateFormField("model")
	if err != nil {
		return "", fmt.Errorf("error creating form field: %w", err)
	}
	fieldWriter.Write([]byte("whisper-1"))

	fieldWriter, err = multiPartWriter.CreateFormField("response_format")
	if err != nil {
		return "", fmt.Errorf("error creating form field: %w", err)
	}
	fieldWriter.Write([]byte("text"))

	// Add the file
	fileWriter, err := multiPartWriter.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", fmt.Errorf("error copying file: %w", err)
	}

	// Close the multipart writer so that the final boundaries are added
	multiPartWriter.Close()

	// Create a client
	client := &http.Client{}

	// Create a new POST request
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set the content type, this is important, the part after the semicolon must be the boundary that was set by the multipart writer
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+multiPartWriter.Boundary())
	// request.Header.Set("Content-Type", "multipart/form-data")

	// Set the other headers
	request.Header.Set("Authorization", "Bearer "+c.openaiToken)

	// Do the request
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error doing request: %w", err)
	}
	defer response.Body.Close()

	// Check the response
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}

// curl --request POST \
//   --url https://api.openai.com/v1/audio/transcriptions \
//   --header "Authorization: Bearer $OPENAI_API_KEY" \
//   --header 'Content-Type: multipart/form-data' \
//   --form file=@/path/to/file/speech.mp3 \
//   --form model=whisper-1 \
//   --form response_format=text
