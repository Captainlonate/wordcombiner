package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Represents one "message" for OpenAI's chat completion API
// Role would be "user" or "system"
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// The struct that is passed to OpenAI during a request.
// This will be marshalled into JSON and sent to the API in the POST body.
type OpenAIRequest struct {
	Model     string          `json:"model"`
	Messages  []OpenAIMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

// Represents the response from OpenAI's chat completion API
// This is the struct that we will unmarshal the response body into.
type ChatCompletion struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// This is the "System" prompt message that will be sent to OpenAI
// on every request. This is very similar to the training
// data used, so that the model will understand the context of the game.
var systemMessage = OpenAIMessage{
	Role:    "system",
	Content: "You are playing a new game called 'Concept Fusion'. The player will provide you with two words or phrases, such as 'Fire' + 'Water'. You will figure out what would be created by combining my two concepts. For example, combining 'Water' and 'Fire' might result in 'Steam 💨'. The combination must always be one word and an emoji, which will be a thing, noun, gerund or adjective. Afterward, you will determine which emoji most closely relates to the new concept that you just generated. If the new concept is 'Steam', then you might choose the emoji '💨'. All of your responses will be one word, followed by a space, followed by exactly one single emoji character. In our previous example, the complete response would be: 'Steam 💨'.",
}

// Makes a network request to the OpenAI chat completion API.
// Passes some messages to the API to set some context and prompt the model.
// Receives a response from the API and unmarshals it into a ChatCompletion struct.
// So, it takes 2 words and returns a ChatCompletion struct.
func GetCombinationFromOpenAI(wordOne string, wordTwo string) (ChatCompletion, error) {
	// Create the request object
	requestPayload := makeOpenAIRequestPayload(wordOne, wordTwo)
	requestObject, err := makeOpenAIRequest(requestPayload)
	if err != nil {
		fmt.Println("Error creating OpenAI request obj:", err)
		return ChatCompletion{}, err
	}

	// Execute the HTTP request to OpenAI
	client := &http.Client{}
	resp, err := client.Do(requestObject)
	if err != nil {
		fmt.Println("Error executing OpenAI request:", err)
		return ChatCompletion{}, err
	}
	defer resp.Body.Close()

	// Read the response body from OpenAI
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading OpenAI response:", err)
		return ChatCompletion{}, err
	}

	// Unmarshal response body into ChatCompletion struct
	var chatCompletion ChatCompletion
	err = json.Unmarshal(responseBodyBytes, &chatCompletion)
	if err != nil {
		fmt.Println("Error unmarshalling OpenAI response body into ChatCompletion struct:", err)
		return ChatCompletion{}, err
	}

	return chatCompletion, nil
}

// Creates the payload that will be sent to OpenAI's chat completion API.
// The only thing that changes in between requests are the two words that the user provides.
func makeOpenAIRequestPayload(wordOne string, wordTwo string) OpenAIRequest {
	return OpenAIRequest{
		Model: os.Getenv("OPENAI_MODEL"),
		Messages: []OpenAIMessage{
			systemMessage,
			{
				Role:    "user",
				Content: fmt.Sprintf("We are playing the 'Concept Fusion' Game. The two words are: '%s' + '%s'. Give me the combined concept and emoji.", wordOne, wordTwo),
			},
		},
		MaxTokens: 300,
	}
}

// To make a request to OpenAI, we need to create an HTTP request object.
// We'll pass this request object to the HTTP client to make the request.
func makeOpenAIRequest(requestPayload OpenAIRequest) (*http.Request, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Marshal the payload into JSON
	requestPayloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling OpenAI payload:", err)
		return nil, err
	}

	// Create a new HTTP request object to send to OpenAI
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewReader(requestPayloadBytes),
	)
	if err != nil {
		fmt.Println("Error creating OpenAI request:", err)
		return nil, err
	}

	// Set the necessary headers before submitting the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req, nil
}
