package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"migrate-history-script/internal/config"
	"migrate-history-script/internal/gemini"
	"migrate-history-script/internal/storage"
	"migrate-history-script/internal/utils"
)

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func testSetup() {
	fmt.Println("🧪 Тестирование настройки миграции истории")

	// Определяем рабочую директорию скрипта
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Ошибка определения рабочей директории: %v", err)
	}

	fmt.Printf("📁 Рабочая директория: %s\n", workDir)

	// Загружаем .env файл
	envPath := filepath.Join(workDir, ".env")
	if err := utils.LoadEnvFile(envPath); err != nil {
		log.Fatalf("❌ Не удалось загрузить .env файл: %v", err)
	}
	fmt.Printf("✅ Конфигурация загружена из: %s\n", envPath)

	// Проверяем обязательные переменные
	requiredVars := []string{
		"POSTGRESQL_HOST",
		"POSTGRESQL_USER",
		"POSTGRESQL_PASSWORD",
		"POSTGRESQL_DBNAME",
		"GEMINI_API_KEY",
	}

	fmt.Println("\n🔧 Проверка переменных окружения:")
	for _, varName := range requiredVars {
		value := os.Getenv(varName)
		if value == "" {
			fmt.Printf("❌ %s: НЕ ЗАДАНА\n", varName)
		} else {
			if varName == "GEMINI_API_KEY" || varName == "POSTGRESQL_PASSWORD" {
				fmt.Printf("✅ %s: ***скрыто*** (длина: %d)\n", varName, len(value))
			} else {
				fmt.Printf("✅ %s: %s\n", varName, value)
			}
		}
	}

	// Проверяем структуру директорий
	fmt.Println("\n📂 Проверка структуры директорий:")

	dirs := map[string]string{
		"data":             filepath.Join(workDir, "data"),
		"cache":            filepath.Join(workDir, "cache"),
		"cache/embeddings": filepath.Join(workDir, "cache", "embeddings"),
	}

	for name, path := range dirs {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("⚠️ %s: %s (НЕ СУЩЕСТВУЕТ)\n", name, path)
			// Создаем директорию
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("❌ Не удалось создать %s: %v\n", name, err)
			} else {
				fmt.Printf("✅ %s: %s (СОЗДАНА)\n", name, path)
			}
		} else {
			fmt.Printf("✅ %s: %s\n", name, path)
		}
	}

	// Проверяем JSON файлы в data
	dataDir := filepath.Join(workDir, "data")
	files, err := filepath.Glob(filepath.Join(dataDir, "*.json"))
	if err != nil {
		fmt.Printf("❌ Ошибка сканирования data директории: %v\n", err)
	} else {
		fmt.Printf("\n📋 JSON файлы в data директории: %d\n", len(files))
		for i, file := range files {
			info, _ := os.Stat(file)
			fmt.Printf("  %d. %s (%s)\n", i+1, filepath.Base(file), formatFileSize(info.Size()))
			if i >= 4 { // Показываем только первые 5 файлов
				fmt.Printf("  ... и еще %d файлов\n", len(files)-5)
				break
			}
		}

		if len(files) == 0 {
			fmt.Println("⚠️ В data директории нет JSON файлов!")
			fmt.Println("   Поместите файлы экспорта Telegram в data/ директорию")
		}
	}

	// Тестируем подключение к PostgreSQL
	fmt.Println("\n🔗 Тестирование подключения к PostgreSQL...")

	dbHost := os.Getenv("POSTGRESQL_HOST")
	dbPort := utils.GetEnvOrDefault("POSTGRESQL_PORT", "5432")
	dbUser := os.Getenv("POSTGRESQL_USER")
	dbPassword := os.Getenv("POSTGRESQL_PASSWORD")
	dbName := os.Getenv("POSTGRESQL_DBNAME")

	storage, err := storage.NewPostgresStorage(dbHost, dbPort, dbUser, dbPassword, dbName, 1000, false)
	if err != nil {
		fmt.Printf("❌ Ошибка подключения к PostgreSQL: %v\n", err)
	} else {
		defer storage.Close()
		fmt.Printf("✅ Успешное подключение к PostgreSQL: %s:%s\n", dbHost, dbPort)
	}

	// Тестируем Gemini API
	fmt.Println("\n🤖 Тестирование Gemini API...")

	cfg := &config.Config{
		GeminiAPIKey:             os.Getenv("GEMINI_API_KEY"),
		GeminiModelName:          "gemini-2.5-flash",
		GeminiEmbeddingModelName: utils.GetEnvOrDefault("GEMINI_EMBEDDING_MODEL_NAME", "embedding-001"),
		Debug:                    false,
	}

	llmClient, err := gemini.New(cfg, cfg.GeminiModelName, cfg.GeminiEmbeddingModelName, cfg.Debug)
	if err != nil {
		fmt.Printf("❌ Ошибка инициализации Gemini: %v\n", err)
	} else {
		defer llmClient.Close()

		// Тестируем генерацию эмбеддинга
		testText := "Это тестовое сообщение для проверки эмбеддинга"
		fmt.Printf("🧪 Тестируем эмбеддинг для: \"%s\"\n", testText)

		start := time.Now()
		embedding, err := llmClient.EmbedContent(testText)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ Ошибка генерации эмбеддинга: %v\n", err)
		} else {
			fmt.Printf("✅ Эмбеддинг сгенерирован успешно:\n")
			fmt.Printf("   • Размерность: %d\n", len(embedding))
			fmt.Printf("   • Время: %v\n", duration)
			fmt.Printf("   • Первые 5 значений: [%.3f, %.3f, %.3f, %.3f, %.3f]\n",
				embedding[0], embedding[1], embedding[2], embedding[3], embedding[4])
		}
	}

	// Показываем настройки производительности
	fmt.Println("\n⚙️ Настройки производительности:")
	fmt.Printf("   • Запросов в минуту: %s\n", utils.GetEnvOrDefault("EMBEDDING_REQUESTS_PER_MINUTE", "240"))
	fmt.Printf("   • Запросов в день: %s\n", utils.GetEnvOrDefault("EMBEDDING_REQUESTS_PER_DAY", "24000"))
	fmt.Printf("   • Задержка между запросами: %s\n", utils.GetEnvOrDefault("EMBEDDING_REQUEST_DELAY", "300ms"))

	// Финальные рекомендации
	fmt.Println("\n📝 Итоги проверки:")

	if len(files) == 0 {
		fmt.Println("⚠️ ВНИМАНИЕ: Добавьте JSON файлы в data/ директорию перед запуском миграции")
	}

	if os.Getenv("GEMINI_API_KEY") == "" {
		fmt.Println("❌ КРИТИЧНО: Задайте GEMINI_API_KEY в .env файле")
	}

	if os.Getenv("POSTGRESQL_HOST") == "" {
		fmt.Println("❌ КРИТИЧНО: Проверьте настройки PostgreSQL в .env файле")
	}

	fmt.Println("\n🚀 Для запуска миграции используйте:")
	fmt.Println("   go run main.go")
	fmt.Println("   или")
	fmt.Println("   go build -o migrate_history.exe && ./migrate_history.exe")
}

func main() {
	testSetup()
}
