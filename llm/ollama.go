package llm

import (
	"context"
	"fmt"
	"strings"

	ollama "github.com/ollama/ollama/api"
)

// OllamaClient wraps the Ollama API client.
type OllamaClient struct {
	client *ollama.Client
}

// NewOllamaClient creates a new OllamaClient.
func NewOllamaClient() (*OllamaClient, error) {
	o, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}
	return &OllamaClient{client: o}, nil
}

// ListModels lists all available Ollama models.
func (oc *OllamaClient) ListModels() ([]string, error) {
	resp, err := oc.client.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}

	var models []string
	for _, model := range resp.Models {
		models = append(models, model.Name)
	}
	return models, nil
}

// GenerateResponse sends a prompt to Ollama and returns the response.
func (oc *OllamaClient) GenerateResponse(modelName, prompt, systemPrompt string) (string, error) {
	var responseText string
	stream := false
	err := oc.client.Generate(context.Background(), &ollama.GenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		System: systemPrompt,
		Stream: &stream,
	}, func(res ollama.GenerateResponse) error {
		responseText += res.Response
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return strings.TrimSpace(responseText), nil
}

// GenerateResponseStream sends a prompt to Ollama and streams the response through a channel.
// It ensures that the channel is closed after the generation is complete.
func (oc *OllamaClient) GenerateResponseStream(modelName, systemPrompt string, history []string, ch chan<- string) {
	defer close(ch)

	messages := []ollama.Message{}
	if systemPrompt != "" {
		messages = append(messages, ollama.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	for i, message := range history {
		var role string
		if i%2 == 0 {
			role = "user"
		} else {
			role = "assistant"
		}
		messages = append(messages, ollama.Message{
			Role:    role,
			Content: message,
		})
	}

	stream := true
	err := oc.client.Chat(context.Background(), &ollama.ChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   &stream,
	}, func(res ollama.ChatResponse) error {
		ch <- res.Message.Content
		return nil
	})

	if err != nil {
		errorMsg := fmt.Sprintf("Error: %v", err)
		ch <- errorMsg
	}
}
