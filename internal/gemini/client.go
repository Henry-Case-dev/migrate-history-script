package gemini

import (
	"context"
	"fmt"

	"migrate-history-script/internal/config"
	"migrate-history-script/internal/llm"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client представляет клиент для работы с Gemini API
type Client struct {
	client         *genai.Client
	embeddingModel *genai.EmbeddingModel
	ctx            context.Context
	debug          bool
}

// New создает новый клиент Gemini
func New(cfg *config.Config, modelName, embeddingModelName string, debug bool) (llm.LLMClient, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	embeddingModel := client.EmbeddingModel(embeddingModelName)

	return &Client{
		client:         client,
		embeddingModel: embeddingModel,
		ctx:            ctx,
		debug:          debug,
	}, nil
}

// EmbedContent генерирует векторное представление для текста
func (c *Client) EmbedContent(text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	embedding, err := c.embeddingModel.EmbedContent(c.ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if embedding == nil || embedding.Embedding == nil {
		return nil, fmt.Errorf("received nil embedding")
	}

	return embedding.Embedding.Values, nil
}

// Close закрывает соединение с клиентом
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
