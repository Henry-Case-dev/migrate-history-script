package llm

// LLMClient предоставляет интерфейс для работы с языковыми моделями
type LLMClient interface {
	// EmbedContent генерирует векторное представление для текста
	EmbedContent(text string) ([]float32, error)

	// Close закрывает соединение с клиентом
	Close() error
}
