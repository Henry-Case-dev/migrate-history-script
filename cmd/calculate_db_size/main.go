package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DBSizeCalculator struct {
	totalMessages   int64
	totalUsers      int64
	avgMessageSize  int64
	avgVectorSize   int64 // 768 float32 values = 768 * 4 = 3072 bytes
	avgUsernameSize int64
	avgTextSize     int64
	totalEmbeddings int64
}

// ChatMessage представляет структуру для расчёта размера сообщения в БД
type ChatMessage struct {
	ChatID              int64     `json:"chat_id"`
	MessageID           int       `json:"message_id"`
	UserID              int64     `json:"user_id"`
	Username            string    `json:"username"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	IsBot               bool      `json:"is_bot"`
	MessageText         string    `json:"message_text"`
	MessageDate         time.Time `json:"message_date"`
	ReplyToMessageID    *int      `json:"reply_to_message_id"`
	Entities            string    `json:"entities"`    // JSON
	RawMessage          string    `json:"raw_message"` // JSON
	IsForward           bool      `json:"is_forward"`
	ForwardedFromUserID *int64    `json:"forwarded_from_user_id"`
	ForwardedFromChatID *int64    `json:"forwarded_from_chat_id"`
	EmbeddingVector     []float32 `json:"message_embedding"` // 768 dimensions
	EmbeddingContext    string    `json:"embedding_context"`
	EmbeddingCreatedAt  time.Time `json:"embedding_generated_at"`
}

func main() {
	fmt.Println("=== Расчёт размера PostgreSQL базы данных для истории чата ===")

	// Читаем состояние миграции
	stateData, err := os.ReadFile("migration_state.json")
	if err != nil {
		log.Fatalf("Ошибка чтения migration_state.json: %v", err)
	}

	var state struct {
		TotalMessages   int64 `json:"total_messages"`
		ProcessedCount  int64 `json:"processed_count"`
		VectorizedCount int64 `json:"vectorized_count"`
	}

	if err := json.Unmarshal(stateData, &state); err != nil {
		log.Fatalf("Ошибка парсинга состояния: %v", err)
	}

	// Анализируем JSON файлы для расчёта среднего размера
	calc := &DBSizeCalculator{
		totalMessages: state.TotalMessages,
		avgVectorSize: 768 * 4, // 768 float32 = 3072 bytes
	}

	fmt.Printf("Общее количество сообщений: %d\n", calc.totalMessages)
	fmt.Printf("Уже обработано: %d\n", state.ProcessedCount)
	fmt.Printf("Уже векторизировано: %d\n", state.VectorizedCount)

	// Анализ одного JSON файла для оценки среднего размера
	err = calc.analyzeJSONFile("../../История чата")
	if err != nil {
		log.Printf("Предупреждение: не удалось проанализировать JSON файлы: %v", err)
		// Используем оценочные значения
		calc.setEstimatedSizes()
	}

	// Расчёт размеров таблиц
	calc.calculateTableSizes()
}

func (calc *DBSizeCalculator) analyzeJSONFile(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("JSON файлы не найдены в %s", dir)
	}

	// Анализируем первый файл для оценки
	fmt.Printf("Анализ файла: %s\n", files[0])
	data, err := os.ReadFile(files[0])
	if err != nil {
		return err
	}

	var chatData struct {
		Messages []tgbotapi.Message `json:"messages"`
	}

	if err := json.Unmarshal(data, &chatData); err != nil {
		return err
	}

	sampleSize := min(1000, len(chatData.Messages)) // Анализируем до 1000 сообщений
	totalTextSize := int64(0)
	totalUsernameSize := int64(0)
	totalRawMessageSize := int64(0)
	uniqueUsers := make(map[int64]bool)

	for i := 0; i < sampleSize; i++ {
		msg := chatData.Messages[i]

		// Размер текста
		totalTextSize += int64(len(msg.Text))

		// Пользователи
		if msg.From != nil {
			uniqueUsers[msg.From.ID] = true
			totalUsernameSize += int64(len(msg.From.UserName) + len(msg.From.FirstName) + len(msg.From.LastName))
		}

		// Размер сырого JSON сообщения
		rawJSON, _ := json.Marshal(msg)
		totalRawMessageSize += int64(len(rawJSON))
	}

	calc.avgTextSize = totalTextSize / int64(sampleSize)
	calc.avgUsernameSize = totalUsernameSize / int64(len(uniqueUsers))
	calc.totalUsers = int64(len(uniqueUsers)) * (calc.totalMessages / int64(sampleSize)) // Экстраполируем

	avgRawMessageSize := totalRawMessageSize / int64(sampleSize)

	fmt.Printf("Анализ образца из %d сообщений:\n", sampleSize)
	fmt.Printf("- Средний размер текста: %d байт\n", calc.avgTextSize)
	fmt.Printf("- Средний размер username+имени: %d байт\n", calc.avgUsernameSize)
	fmt.Printf("- Средний размер raw_message JSON: %d байт\n", avgRawMessageSize)
	fmt.Printf("- Уникальных пользователей в образце: %d\n", len(uniqueUsers))
	fmt.Printf("- Оценка общего количества пользователей: %d\n", calc.totalUsers)

	return nil
}

func (calc *DBSizeCalculator) setEstimatedSizes() {
	// Оценочные значения на основе типичных Telegram чатов
	calc.avgTextSize = 150                    // байт (средняя длина сообщения)
	calc.avgUsernameSize = 50                 // байт (username + имя + фамилия)
	calc.totalUsers = calc.totalMessages / 50 // примерно 50 сообщений на пользователя

	fmt.Println("Используются оценочные значения:")
	fmt.Printf("- Средний размер текста: %d байт\n", calc.avgTextSize)
	fmt.Printf("- Средний размер имён пользователей: %d байт\n", calc.avgUsernameSize)
	fmt.Printf("- Оценка количества пользователей: %d\n", calc.totalUsers)
}

func (calc *DBSizeCalculator) calculateTableSizes() {
	fmt.Println("\n=== РАСЧЁТ РАЗМЕРА ТАБЛИЦ ===")

	// Таблица chat_messages
	msgTableSize := calc.calculateChatMessagesSize()

	// Таблица user_profiles
	userTableSize := calc.calculateUserProfilesSize()

	// Индексы
	indexesSize := calc.calculateIndexesSize()

	totalSize := msgTableSize + userTableSize + indexesSize

	fmt.Printf("\n=== ИТОГОВЫЙ РАЗМЕР ===\n")
	fmt.Printf("Таблица chat_messages: %s\n", formatBytes(msgTableSize))
	fmt.Printf("Таблица user_profiles: %s\n", formatBytes(userTableSize))
	fmt.Printf("Индексы: %s\n", formatBytes(indexesSize))
	fmt.Printf("ОБЩИЙ РАЗМЕР БД: %s\n", formatBytes(totalSize))

	// Дополнительная информация
	fmt.Printf("\n=== ДОПОЛНИТЕЛЬНАЯ ИНФОРМАЦИЯ ===\n")
	fmt.Printf("Размер векторных эмбеддингов: %s\n", formatBytes(calc.totalMessages*calc.avgVectorSize))
	fmt.Printf("Процент от общего размера (векторы): %.1f%%\n",
		float64(calc.totalMessages*calc.avgVectorSize)/float64(totalSize)*100)

	// Сравнение с исходными JSON файлами
	estimatedJSONSize := calc.totalMessages * 500 // примерно 500 байт на сообщение в JSON
	fmt.Printf("Размер исходных JSON (оценка): %s\n", formatBytes(estimatedJSONSize))
	fmt.Printf("Отношение БД/JSON: %.1fx\n", float64(totalSize)/float64(estimatedJSONSize))
}

func (calc *DBSizeCalculator) calculateChatMessagesSize() int64 {
	fmt.Println("\n--- Таблица chat_messages ---")

	// Фиксированные поля
	fixedSize := int64(8 + 4 + 8 + 8 + 1 + 8 + 4 + 1 + 8 + 8 + 4 + 8) // 70 байт основных полей

	// Переменные поля
	avgUsernameFields := calc.avgUsernameSize * 3 // username, first_name, last_name
	avgTextSize := calc.avgTextSize
	avgEntitiesSize := int64(100)     // JSON entities
	avgRawMessageSize := int64(800)   // полное JSON сообщение
	avgEmbeddingContext := int64(500) // контекст для эмбеддинга
	vectorSize := calc.avgVectorSize  // 768 float32

	rowSize := fixedSize + avgUsernameFields + avgTextSize + avgEntitiesSize +
		avgRawMessageSize + avgEmbeddingContext + vectorSize

	totalSize := calc.totalMessages * rowSize

	// Overhead PostgreSQL (примерно 15-20%)
	overhead := int64(float64(totalSize) * 0.18)
	totalWithOverhead := totalSize + overhead

	fmt.Printf("Размер строки (средний): %d байт\n", rowSize)
	fmt.Printf("  - Фиксированные поля: %d байт\n", fixedSize)
	fmt.Printf("  - Username/имена: %d байт\n", avgUsernameFields)
	fmt.Printf("  - Текст сообщения: %d байт\n", avgTextSize)
	fmt.Printf("  - Entities JSON: %d байт\n", avgEntitiesSize)
	fmt.Printf("  - Raw message JSON: %d байт\n", avgRawMessageSize)
	fmt.Printf("  - Embedding context: %d байт\n", avgEmbeddingContext)
	fmt.Printf("  - Vector (768 float32): %d байт\n", vectorSize)
	fmt.Printf("Чистый размер данных: %s\n", formatBytes(totalSize))
	fmt.Printf("PostgreSQL overhead (~18%%): %s\n", formatBytes(overhead))
	fmt.Printf("Общий размер таблицы: %s\n", formatBytes(totalWithOverhead))

	return totalWithOverhead
}

func (calc *DBSizeCalculator) calculateUserProfilesSize() int64 {
	fmt.Println("\n--- Таблица user_profiles ---")

	// Поля user_profiles
	fixedSize := int64(8 + 8 + 8 + 8 + 8)                              // chat_id, user_id, timestamps
	avgFieldsSize := calc.avgUsernameSize + 50 + 50 + 200 + 500 + 1000 // username, alias, gender, real_name, bio, auto_bio

	rowSize := fixedSize + avgFieldsSize
	totalSize := calc.totalUsers * rowSize

	// PostgreSQL overhead
	overhead := int64(float64(totalSize) * 0.15)
	totalWithOverhead := totalSize + overhead

	fmt.Printf("Пользователей: %d\n", calc.totalUsers)
	fmt.Printf("Размер строки (средний): %d байт\n", rowSize)
	fmt.Printf("Чистый размер данных: %s\n", formatBytes(totalSize))
	fmt.Printf("PostgreSQL overhead (~15%%): %s\n", formatBytes(overhead))
	fmt.Printf("Общий размер таблицы: %s\n", formatBytes(totalWithOverhead))

	return totalWithOverhead
}

func (calc *DBSizeCalculator) calculateIndexesSize() int64 {
	fmt.Println("\n--- Индексы ---")

	// Первичные ключи и основные индексы
	primaryKeySize := calc.totalMessages * 12 // (chat_id + message_id)
	dateIndexSize := calc.totalMessages * 16  // chat_id + timestamp
	userIndexSize := calc.totalMessages * 8   // user_id
	profileIndexSize := calc.totalUsers * 16  // chat_id + user_id

	// Векторный HNSW индекс (самый большой)
	// HNSW индекс обычно занимает 1.5-2x размера векторных данных
	vectorDataSize := calc.totalMessages * calc.avgVectorSize
	hsnwIndexSize := int64(float64(vectorDataSize) * 1.8)

	totalIndexSize := primaryKeySize + dateIndexSize + userIndexSize + profileIndexSize + hsnwIndexSize

	fmt.Printf("Primary keys: %s\n", formatBytes(primaryKeySize))
	fmt.Printf("Date index: %s\n", formatBytes(dateIndexSize))
	fmt.Printf("User index: %s\n", formatBytes(userIndexSize))
	fmt.Printf("Profile index: %s\n", formatBytes(profileIndexSize))
	fmt.Printf("HNSW vector index: %s\n", formatBytes(hsnwIndexSize))
	fmt.Printf("Общий размер индексов: %s\n", formatBytes(totalIndexSize))

	return totalIndexSize
}

func formatBytes(bytes int64) string {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
