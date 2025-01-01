package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Request structure
type ChatGPTRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// Response structure
type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("API")
	// using free openai plan so the api quota is limited.
	if apiKey == "" {
		log.Fatal("API key not found. Please set it in the .env file as 'API'.")
	}

	url := "https://api.openai.com/v1/chat/completions"

	// Prompt user for a query
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your query: ")
	userQuery, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	// request payload
	requestBody := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
	}
	requestBody.Messages = append(requestBody.Messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: userQuery,
	})

	// Convert to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Failed to create request body: %v", err)
	}

	// Make the API call
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read and response processing
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var chatResponse ChatGPTResponse
	err = json.Unmarshal(responseBody, &chatResponse)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	// response
	if len(chatResponse.Choices) > 0 {
		fmt.Println("ChatGPT Response:", chatResponse.Choices[0].Message.Content)
	} else {
		fmt.Println("No response from ChatGPT.")
	}
}
