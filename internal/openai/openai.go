package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
	}
}

func (s *OpenAIService) GetAIResponse(history []string) (string, error) {
	message := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "Ты полезный AI-бот.",
		},
	}
	for _, msg := range history {
		message = append(message, openai.ChatCompletionMessage{
			Role:    "user",
			Content: msg,
		})
	}
	response, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: message,
	},
	)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
