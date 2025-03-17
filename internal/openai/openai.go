package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type OpenAIService struct {
	apiKey   string
	proxyURL string
}

func NewOpenAIService(apiKey, proxyURL string) *OpenAIService {
	return &OpenAIService{
		apiKey:   apiKey,
		proxyURL: proxyURL,
	}
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageGenerationResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (s *OpenAIService) GetAIResponse(history []string) (string, error) {
	messages := []Message{
		{Role: "system", Content: "Ты полезный AI-бот."},
	}

	for _, msg := range history {
		messages = append(messages, Message{Role: "user", Content: msg})
	}

	reqBody := ChatCompletionRequest{
		Model:    "gpt-4o", // model
		Messages: messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Отправка запроса через ProxyAPI
	req, err := http.NewRequest("POST", s.proxyURL+"/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка запроса: статус %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("не получен ответ от ProxyAPI")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("неверная структура ответа")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", errors.New("не найдено поле 'message'")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("не найдено содержимое сообщения")
	}

	return content, nil
}

func (s *OpenAIService) GenerateImage(prompt string) ([]string, error) {
	reqBody := ImageGenerationRequest{
		Model:  "dall-e-3", // model
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.proxyURL+"/openai/v1/images/generations", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка запроса: статус %d", resp.StatusCode)
	}

	var response ImageGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, errors.New("не получен ответ от ProxyAPI")
	}

	return []string{response.Data[0].URL}, nil
}
