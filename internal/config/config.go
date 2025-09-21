package config

// Config содержит конфигурацию для работы со сторонними сервисами
type Config struct {
	// Gemini API
	GeminiAPIKey             string
	GeminiModelName          string
	GeminiEmbeddingModelName string
	
	// PostgreSQL
	PostgreSQLHost     string
	PostgreSQLPort     string
	PostgreSQLUser     string
	PostgreSQLPassword string
	PostgreSQLDBName   string
	
	// Общие настройки
	Debug bool
}