package openai

import (
	"MovieTrackerBack/internal/domain"
	"MovieTrackerBack/internal/infrastructure/postgres"
	"MovieTrackerBack/internal/infrastructure/redis"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/openai/openai-go"
	openai_go "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client         *openai.Client
	redis          *redis.CacheRepository
	model          string
	itemRepository *postgres.ItemRepository
	listRepository *postgres.ListRepository
}

func NewOpenAIClient(redis *redis.CacheRepository, api_key string, model string, itemRepository *postgres.ItemRepository, listRepository *postgres.ListRepository) *OpenAIClient {
	openAIGoClient := openai_go.NewClient(option.WithAPIKey(api_key))
	return &OpenAIClient{client: &openAIGoClient, redis: redis, model: model, itemRepository: itemRepository, listRepository: listRepository}
}

func (c *OpenAIClient) GenerateMessage(msg domain.ChatMessage) (string, error) {
	ctx := context.Background()
	var history []domain.ChatMessage
	data, err := c.redis.Get("chat_" + msg.SessionID)
	if err == nil && data != nil {
		json.Unmarshal(data, &history)
	}

	var messages []openai.ChatCompletionMessageParamUnion
	messages = append(messages, openai.SystemMessage(
		"Eres un asistente personal de películas y series.",
	))

	for _, h := range history {
		if h.Role == "user" {
			messages = append(messages, openai.UserMessage(h.Content))
		} else {
			messages = append(messages, openai.AssistantMessage(h.Content))
		}
	}

	messages = append(messages, openai.UserMessage(msg.Content))
	response, err := c.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(c.model),
			Messages: messages,
			Tools:    c.getTools(),
		},
	)
	if err != nil {
		return "", err
	}

	respContent := response.Choices[0].Message.Content

	history = append(history, msg)
	history = append(history, domain.ChatMessage{
		Role:    "assistant",
		Content: respContent,
	})

	newData, _ := json.Marshal(history)
	c.redis.Set("chat_"+msg.SessionID, newData, 2*time.Hour)

	return respContent, nil
}

func (c *OpenAIClient) GetHistory(sessionID string) ([]domain.ChatMessage, error) {
	// TODO: Implementar la obtención de historial desde Redis u otra DB
	return []domain.ChatMessage{}, nil
}

// Tools
func (c *OpenAIClient) getTools() []openai.ChatCompletionToolParam {
	return []openai.ChatCompletionToolParam{
		{
			Type: "function",
			Function: openai.FunctionDefinitionParam{
				Name:        "obtener_watchlist",
				Description: openai.String("Obtiene toda la watchlist del usuario, tanto películas como series, vistas y pendientes"),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionDefinitionParam{
				Name:        "obtener_items_vistos",
				Description: openai.String("Obtiene solo las películas y series que el usuario ya vio, con su calificación y comentario personal"),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionDefinitionParam{
				Name:        "obtener_items_pendientes",
				Description: openai.String("Obtiene solo las películas y series que el usuario tiene pendientes de ver"),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
	}
}

// Ejecutar tools
func (c *OpenAIClient) executeTool(name string) (string, error) {
	switch name {
	case "obtener_watchlist":
		items, err := c.itemRepository.GetAll(domain.Filters{}))
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(items)
		return string(data), nil

	case "obtener_items_vistos":
		items, err := c.getItems("watched")
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(items)
		return string(data), nil

	case "obtener_items_pendientes":
		items, err := c.getItems("pending")
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(items)
		return string(data), nil
	}

	return "", errors.New("tool no encontrada: " + name)
}
