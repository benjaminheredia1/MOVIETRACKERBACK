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
	
	systemPrompt := `Eres un asistente personal experto en películas y series, integrado en una aplicación estilo Letterboxd/Trakt.
Tienes acceso directo a la base de datos personal del usuario mediante herramientas (tools).
Tu trabajo es responder a las consultas del usuario usando esta información y tu propio conocimiento sobre cine y televisión.

Reglas:
- Si te preguntan por las películas o series que el usuario ha guardado (vistas o pendientes), usa "obtener_items".
- Si te preguntan por las listas personalizadas que el usuario ha creado, usa "obtener_listas".
- Siempre responde de forma amigable, útil, y en formato Markdown.
- Da recomendaciones basadas en los géneros o películas que ya han visto si te piden sugerencias.`

	messages = append(messages, openai.SystemMessage(systemPrompt))

	for _, h := range history {
		if h.Role == "user" {
			messages = append(messages, openai.UserMessage(h.Content))
		} else {
			messages = append(messages, openai.AssistantMessage(h.Content))
		}
	}

	messages = append(messages, openai.UserMessage(msg.Content))

	for {
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

		choice := response.Choices[0]
		
		// Añadir el mensaje del asistente (que puede contener la petición de tool call) al historial
		messages = append(messages, choice.Message.ToParam())

		if len(choice.Message.ToolCalls) > 0 {
			for _, toolCall := range choice.Message.ToolCalls {
				toolResult, err := c.executeTool(toolCall.Function.Name)
				if err != nil {
					toolResult = "Error al ejecutar la herramienta: " + err.Error()
				}
				messages = append(messages, openai.ToolMessage(toolCall.ID, toolResult))
			}
			continue // Volver a llamar a OpenAI con los resultados
		}

		respContent := choice.Message.Content

		history = append(history, msg)
		history = append(history, domain.ChatMessage{
			Role:    "assistant",
			Content: respContent,
		})

		newData, _ := json.Marshal(history)
		c.redis.Set("chat_"+msg.SessionID, newData, 2*time.Hour)

		return respContent, nil
	}
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
				Name:        "obtener_items",
				Description: openai.String("Obtiene las películas y series del usuario (vistas y pendientes) con sus calificaciones"),
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
				Name:        "obtener_listas",
				Description: openai.String("Obtiene las listas personalizadas del usuario"),
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
	case "obtener_items":
		items, err := c.itemRepository.GetAll(domain.Filters{})
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(items)
		return string(data), nil

	case "obtener_listas":
		listas, err := c.listRepository.GetAll(domain.Filters{})
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(listas)
		return string(data), nil
	}

	return "", errors.New("tool no encontrada: " + name)
}
