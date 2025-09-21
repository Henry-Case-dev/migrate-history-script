package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TelegramMessage упрощённая структура для парсинга JSON истории
type TelegramMessage struct {
	ID               int         `json:"id"`
	Type             string      `json:"type"`
	Date             string      `json:"date"`
	DateUnixtime     int         `json:"date_unixtime"`
	From             string      `json:"from"`
	FromID           string      `json:"from_id"`
	Text             interface{} `json:"text"` // может быть строкой или массивом объектов
	ReplyToMessageID int         `json:"reply_to_message_id"`
}

type ChatHistory struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	ID       int64             `json:"id"`
	Messages []TelegramMessage `json:"messages"`
}

func main() {
	fmt.Println("=== ДЕТАЛЬНЫЙ АНАЛИЗ РАЗМЕРА БД ===")

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

	// Анализируем JSON файлы
	totalMessages := state.TotalMessages
	fmt.Printf("Общее количество сообщений: %d\n", totalMessages)

	// Анализ размера JSON файлов
	analyzeJSONFiles()

	// Детальные расчёты для БД
	calculateDatabaseSizes(totalMessages)
}

func analyzeJSONFiles() {
	fmt.Println("\n=== АНАЛИЗ ИСХОДНЫХ JSON ФАЙЛОВ ===")

	files, err := filepath.Glob("../../История чата/*.json")
	if err != nil {
		log.Printf("Ошибка поиска JSON файлов: %v", err)
		return
	}

	totalSize := int64(0)
	totalFiles := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		totalSize += info.Size()
		totalFiles++
		fmt.Printf("- %s: %s\n", filepath.Base(file), formatBytes(info.Size()))
	}

	fmt.Printf("\nВсего JSON файлов: %d\n", totalFiles)
	fmt.Printf("Общий размер JSON: %s\n", formatBytes(totalSize))

	if totalFiles > 0 {
		avgFileSize := totalSize / int64(totalFiles)
		fmt.Printf("Средний размер файла: %s\n", formatBytes(avgFileSize))
	}
}

func calculateDatabaseSizes(totalMessages int64) {
	fmt.Println("\n=== РАСЧЁТ РАЗМЕРА POSTGRESQL БД ===")

	// Более точные оценки на основе анализа Telegram JSON формата
	avgTextSize := int64(200)          // текст сообщения
	avgUsernameSize := int64(40)       // имя пользователя
	avgEntitiesSize := int64(50)       // entities (обычно небольшие)
	avgRawMessageSize := int64(600)    // полный JSON сообщения
	vectorSize := int64(768 * 4)       // 768 float32 = 3072 байта
	embeddingContextSize := int64(400) // контекст для векторизации

	estimatedUsers := totalMessages / 45 // ~45 сообщений на пользователя (типичная статистика)

	fmt.Printf("Оценка количества пользователей: %d\n", estimatedUsers)
	fmt.Printf("Среднее количество сообщений на пользователя: %.1f\n", float64(totalMessages)/float64(estimatedUsers))

	// 1. Таблица chat_messages
	fmt.Println("\n--- Таблица chat_messages ---")

	// Структура строки в PostgreSQL:
	// - Фиксированные поля: chat_id(8) + message_id(4) + user_id(8) + is_bot(1) + message_date(8) + reply_to_message_id(4) + is_forward(1) + другие(20) = 54 байта
	fixedFields := int64(54)

	// Переменные поля
	variableFields := avgUsernameSize*3 + avgTextSize + avgEntitiesSize + avgRawMessageSize + embeddingContextSize + vectorSize

	rowSize := fixedFields + variableFields
	dataSize := totalMessages * rowSize

	// PostgreSQL overhead (page header, tuple header, alignment, vacuum space)
	overhead := int64(float64(dataSize) * 0.20) // 20% overhead
	tableSize := dataSize + overhead

	fmt.Printf("Размер строки: %d байт\n", rowSize)
	fmt.Printf("  Фиксированные поля: %d байт\n", fixedFields)
	fmt.Printf("  Username/имена (3 поля): %d байт\n", avgUsernameSize*3)
	fmt.Printf("  Текст сообщения: %d байт\n", avgTextSize)
	fmt.Printf("  Entities JSON: %d байт\n", avgEntitiesSize)
	fmt.Printf("  Raw message JSON: %d байт\n", avgRawMessageSize)
	fmt.Printf("  Embedding context: %d байт\n", embeddingContextSize)
	fmt.Printf("  Vector embedding: %d байт\n", vectorSize)
	fmt.Printf("Чистые данные: %s\n", formatBytes(dataSize))
	fmt.Printf("PostgreSQL overhead (20%%): %s\n", formatBytes(overhead))
	fmt.Printf("Размер таблицы: %s\n", formatBytes(tableSize))

	// 2. Таблица user_profiles
	fmt.Println("\n--- Таблица user_profiles ---")

	profileFixedFields := int64(40)                                // chat_id + user_id + timestamps
	profileVariableFields := int64(40 + 50 + 20 + 100 + 300 + 800) // username + alias + gender + real_name + bio + auto_bio
	profileRowSize := profileFixedFields + profileVariableFields
	profileDataSize := estimatedUsers * profileRowSize
	profileOverhead := int64(float64(profileDataSize) * 0.15)
	profileTableSize := profileDataSize + profileOverhead

	fmt.Printf("Пользователей: %d\n", estimatedUsers)
	fmt.Printf("Размер строки профиля: %d байт\n", profileRowSize)
	fmt.Printf("Чистые данные: %s\n", formatBytes(profileDataSize))
	fmt.Printf("PostgreSQL overhead (15%%): %s\n", formatBytes(profileOverhead))
	fmt.Printf("Размер таблицы: %s\n", formatBytes(profileTableSize))

	// 3. Индексы
	fmt.Println("\n--- Индексы ---")

	// B-tree индексы
	primaryKeyIndex := totalMessages * 12      // (chat_id, message_id)
	dateIndex := totalMessages * 16            // (chat_id, message_date)
	userIndex := totalMessages * 8             // user_id
	profilePrimaryIndex := estimatedUsers * 16 // (chat_id, user_id)
	profileUniqueIndex := estimatedUsers * 16  // unique constraint

	btreeIndexesSize := primaryKeyIndex + dateIndex + userIndex + profilePrimaryIndex + profileUniqueIndex

	// HNSW векторный индекс (самый объёмный)
	vectorDataSize := totalMessages * vectorSize
	// HNSW индекс обычно занимает 1.5-2.5x от размера векторных данных, используем 2.0x
	hsnwIndexSize := vectorDataSize * 2

	totalIndexSize := btreeIndexesSize + hsnwIndexSize

	fmt.Printf("B-tree индексы: %s\n", formatBytes(btreeIndexesSize))
	fmt.Printf("  - Primary key (messages): %s\n", formatBytes(primaryKeyIndex))
	fmt.Printf("  - Date index: %s\n", formatBytes(dateIndex))
	fmt.Printf("  - User index: %s\n", formatBytes(userIndex))
	fmt.Printf("  - Profile indexes: %s\n", formatBytes(profilePrimaryIndex+profileUniqueIndex))
	fmt.Printf("HNSW vector index: %s\n", formatBytes(hsnwIndexSize))
	fmt.Printf("Общий размер индексов: %s\n", formatBytes(totalIndexSize))

	// Общий размер
	totalDBSize := tableSize + profileTableSize + totalIndexSize

	fmt.Println("\n=== ИТОГОВЫЙ РАЗМЕР БД ===")
	fmt.Printf("Таблица chat_messages: %s\n", formatBytes(tableSize))
	fmt.Printf("Таблица user_profiles: %s\n", formatBytes(profileTableSize))
	fmt.Printf("Индексы: %s\n", formatBytes(totalIndexSize))
	fmt.Printf("ОБЩИЙ РАЗМЕР: %s\n", formatBytes(totalDBSize))

	// Дополнительная аналитика
	fmt.Println("\n=== АНАЛИТИКА ===")
	fmt.Printf("Размер только векторных данных: %s\n", formatBytes(vectorDataSize))
	fmt.Printf("Доля векторов в общем размере: %.1f%%\n", float64(vectorDataSize+hsnwIndexSize)/float64(totalDBSize)*100)

	// Сравнение с JSON
	estimatedJSONSize := totalMessages * 500 // примерная оценка
	fmt.Printf("Размер исходных JSON (оценка): %s\n", formatBytes(estimatedJSONSize))
	fmt.Printf("Коэффициент увеличения БД/JSON: %.1fx\n", float64(totalDBSize)/float64(estimatedJSONSize))

	// Дополнительные метрики
	fmt.Printf("Размер БД на одно сообщение: %.1f KB\n", float64(totalDBSize)/float64(totalMessages)/1024)
	fmt.Printf("Размер БД на одного пользователя: %.1f MB\n", float64(totalDBSize)/float64(estimatedUsers)/1024/1024)

	// Требования к диску с запасом
	fmt.Printf("\nРЕКОМЕНДУЕМЫЙ РАЗМЕР ДИСКА: %s (с запасом 50%%)\n", formatBytes(int64(float64(totalDBSize)*1.5)))
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
