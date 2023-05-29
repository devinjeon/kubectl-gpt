package gpt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenAIRequest(model string, temperature float64, maxTokens int, systemMessage, userMessage string) OpenAIRequest {
	return OpenAIRequest{
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemMessage,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}
}

func RequestChatGptAPI(reqBody OpenAIRequest, apiKey string) (OpenAIResponse, error) {
	client := &http.Client{}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil { return OpenAIResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(jsonReqBody)))
	if err != nil {
		return OpenAIResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return OpenAIResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return OpenAIResponse{}, fmt.Errorf("received status code %d, body: %s", resp.StatusCode, bodyString)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OpenAIResponse{}, err
	}

	var respData OpenAIResponse
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return OpenAIResponse{}, err
	}

	return respData, nil
}
