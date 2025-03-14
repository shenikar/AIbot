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

func (s *OpenAIService) GetAIResponse(history []string) (string, error) {
	messages := []Message{
		{Role: "system", Content: "Ты полезный AI-бот."},
	}

	for _, msg := range history {
		messages = append(messages, Message{Role: "user", Content: msg})
	}

	reqBody := ChatCompletionRequest{
		Model:    "gpt-3.5-turbo", // Используйте модель, подходящую для ProxyAPI
		Messages: messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Отправка запроса через ProxyAPI
	req, err := http.NewRequest("POST", s.proxyURL, bytes.NewBuffer(body))
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
