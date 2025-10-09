package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"migrate-history-script/internal/config"
	"migrate-history-script/internal/gemini"
	"migrate-history-script/internal/llm"
	"migrate-history-script/internal/storage"
	"migrate-history-script/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/schollz/progressbar/v3"
)

// Структуры для парсинга экспорта Telegram
type TelegramExportMessage struct {
	ID               int           `json:"id"`
	Type             string        `json:"type"`
	Date             string        `json:"date"`
	DateUnixtime     string        `json:"date_unixtime"`
	From             string        `json:"from,omitempty"`
	FromID           string        `json:"from_id,omitempty"`
	Text             interface{}   `json:"text,omitempty"`
	ReplyToMessageID int           `json:"reply_to_message_id,omitempty"`
	Photo            string        `json:"photo,omitempty"`
	File             string        `json:"file,omitempty"`
	MediaType        string        `json:"media_type,omitempty"`
	Entities         []interface{} `json:"entities,omitempty"`
}

type TelegramExport struct {
	Name     string                  `json:"name"`
	Type     string                  `json:"type"`
	ID       int64                   `json:"id"`
	Messages []TelegramExportMessage `json:"messages"`
}

// Состояние миграции для возобновления
type MigrationState struct {
	CurrentFile     string        `json:"current_file"`
	ProcessedFiles  []string      `json:"processed_files"`
	CurrentIndex    int           `json:"current_index"`
	TotalMessages   int           `json:"total_messages"`
	ProcessedCount  int           `json:"processed_count"`
	VectorizedCount int           `json:"vectorized_count"`
	LastUpdateTime  time.Time     `json:"last_update_time"`
	ChatMessagesMap map[int64]int `json:"chat_messages_map"` // chatID -> количество сообщений
}

// Rate Limiter для API запросов
type RateLimiter struct {
	requestsThisMinute   int
	requestsToday        int
	lastMinuteReset      time.Time
	lastDayReset         time.Time
	maxRequestsPerMinute int
	maxRequestsPerDay    int
	requestDelay         time.Duration
	currentDelay         time.Duration
	consecutiveErrors    int
	isThrottled          bool
	mu                   sync.Mutex
}

func NewRateLimiter(reqPerMin, reqPerDay int, delay time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequestsPerMinute: reqPerMin,
		maxRequestsPerDay:    reqPerDay,
		requestDelay:         delay,
		currentDelay:         delay,
		lastMinuteReset:      time.Now(),
		lastDayReset:         time.Now(),
	}
}

func (rl *RateLimiter) WaitForRequest(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Сброс счетчиков
	if now.Sub(rl.lastMinuteReset) >= time.Minute {
		rl.requestsThisMinute = 0
		rl.lastMinuteReset = now
	}

	if now.Sub(rl.lastDayReset) >= 24*time.Hour {
		rl.requestsToday = 0
		rl.lastDayReset = now
	}

	// Проверка лимитов
	if rl.requestsThisMinute >= rl.maxRequestsPerMinute {
		waitTime := time.Minute - now.Sub(rl.lastMinuteReset)
		log.Printf("[RateLimiter] Достигнут минутный лимит, ожидание %v", waitTime)

		select {
		case <-time.After(waitTime):
			rl.requestsThisMinute = 0
			rl.lastMinuteReset = time.Now()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if rl.requestsToday >= rl.maxRequestsPerDay {
		waitTime := 24*time.Hour - now.Sub(rl.lastDayReset)
		log.Printf("[RateLimiter] Достигнут дневной лимит, ожидание %v", waitTime)

		select {
		case <-time.After(waitTime):
			rl.requestsToday = 0
			rl.lastDayReset = time.Now()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Адаптивная задержка
	if rl.isThrottled {
		select {
		case <-time.After(rl.currentDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Увеличиваем счетчики
	rl.requestsThisMinute++
	rl.requestsToday++

	return nil
}

func (rl *RateLimiter) OnSuccess() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.consecutiveErrors = 0
	rl.isThrottled = false
	rl.currentDelay = rl.requestDelay
}

func (rl *RateLimiter) OnRateLimit() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.consecutiveErrors++
	rl.isThrottled = true

	// Экспоненциальное увеличение задержки
	rl.currentDelay = time.Duration(float64(rl.requestDelay) * (1.5 * float64(rl.consecutiveErrors)))
	if rl.currentDelay > 5*time.Minute {
		rl.currentDelay = 5 * time.Minute
	}

	log.Printf("[RateLimiter] Rate limit detected, увеличена задержка до %v", rl.currentDelay)
}

func (rl *RateLimiter) GetStatus() (int, int, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return rl.requestsThisMinute, rl.requestsToday, rl.currentDelay
}

// Контекстный генератор эмбеддингов
type ContextualEmbedder struct {
	llmClient   llm.LLMClient
	rateLimiter *RateLimiter
	cacheDir    string
	cache       map[string][]float32
	cacheMu     sync.RWMutex
}

func NewContextualEmbedder(llmClient llm.LLMClient, rateLimiter *RateLimiter, cacheDir string) *ContextualEmbedder {
	ce := &ContextualEmbedder{
		llmClient:   llmClient,
		rateLimiter: rateLimiter,
		cacheDir:    cacheDir,
		cache:       make(map[string][]float32),
	}

	// Создаем директорию для кэша
	if cacheDir != "" {
		os.MkdirAll(cacheDir, 0755)
		ce.loadCache()
	}

	return ce
}

func (ce *ContextualEmbedder) GenerateEmbeddingWithContext(
	currentMsg *TelegramExportMessage,
	contextMsgs []TelegramExportMessage,
	targetIndex int,
) ([]float32, string, error) {

	// Формируем контекстную строку
	contextStr := ce.buildContextString(contextMsgs, targetIndex)

	// Проверяем кэш
	if ce.cacheDir != "" {
		if cached := ce.getCachedEmbedding(contextStr); cached != nil {
			return cached, contextStr, nil
		}
	}

	// Ждем разрешения на запрос
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := ce.rateLimiter.WaitForRequest(ctx); err != nil {
		return nil, "", fmt.Errorf("timeout waiting for rate limiter: %w", err)
	}

	// Генерируем эмбеддинг
	embedding, err := ce.llmClient.EmbedContent(contextStr)
	if err != nil {
		ce.rateLimiter.OnRateLimit()
		return nil, "", fmt.Errorf("ошибка генерации эмбеддинга: %w", err)
	}

	ce.rateLimiter.OnSuccess()

	// Сохраняем в кэш
	if ce.cacheDir != "" {
		ce.setCachedEmbedding(contextStr, embedding)
	}

	return embedding, contextStr, nil
}

func (ce *ContextualEmbedder) buildContextString(msgs []TelegramExportMessage, targetIndex int) string {
	if len(msgs) == 0 || targetIndex >= len(msgs) {
		return ""
	}

	var builder strings.Builder
	windowSize := ce.determineOptimalWindowSize(&msgs[targetIndex], msgs, targetIndex)

	start := max(0, targetIndex-windowSize)
	end := min(len(msgs), targetIndex+windowSize+1)

	for i := start; i < end; i++ {
		msg := &msgs[i]
		if !ce.shouldIncludeInContext(msg) {
			continue
		}

		userName := ce.getUserDisplayName(msg)
		text := ce.getMessageText(msg)

		if i == targetIndex {
			// Выделяем целевое сообщение
			builder.WriteString(fmt.Sprintf(">>> %s: %s <<<\n", userName, text))
		} else {
			builder.WriteString(fmt.Sprintf("%s: %s\n", userName, text))
		}
	}

	return strings.TrimSpace(builder.String())
}

func (ce *ContextualEmbedder) determineOptimalWindowSize(msg *TelegramExportMessage, context []TelegramExportMessage, index int) int {
	baseWindow := 5

	// Увеличиваем окно для важных сообщений
	if ce.isImportantMessage(msg) {
		return 25
	}

	// Увеличиваем окно при наличии reply chains
	if msg.ReplyToMessageID != 0 {
		return 15
	}

	text := ce.getMessageText(msg)
	if len(text) > 100 || strings.Contains(text, "?") {
		return 10
	}

	return baseWindow
}

func (ce *ContextualEmbedder) isImportantMessage(msg *TelegramExportMessage) bool {
	text := ce.getMessageText(msg)
	if len(text) > 200 {
		return true
	}

	// Сообщения с вопросами
	if strings.Contains(text, "?") || strings.Contains(text, "как") || strings.Contains(text, "что") {
		return true
	}

	// Сообщения с ключевыми словами
	keywords := []string{"важно", "проблема", "помощь", "решение", "задача", "пиздец", "норм", "кринж", "кайф", "лол", "зачем", "почему", "кек"}
	lowerText := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}

	return false
}

func (ce *ContextualEmbedder) shouldIncludeInContext(msg *TelegramExportMessage) bool {
	if msg.Type != "message" {
		return false
	}

	text := strings.TrimSpace(ce.getMessageText(msg))
	if len(text) < 5 {
		return false
	}

	// Исключаем повторяющиеся паттерны
	lowText := strings.ToLower(text)
	if lowText == "да" || lowText == "нет" || lowText == "ок" || lowText == "лол" || lowText == "👍" || lowText == "👎" {
		return false
	}

	return true
}

func (ce *ContextualEmbedder) getUserDisplayName(msg *TelegramExportMessage) string {
	if msg.From != "" {
		return msg.From
	}
	if msg.FromID != "" {
		return msg.FromID
	}
	return "Unknown"
}

func (ce *ContextualEmbedder) getMessageText(msg *TelegramExportMessage) string {
	if msg.Text == nil {
		return "[non-text message]"
	}

	switch v := msg.Text.(type) {
	case string:
		return v
	case []interface{}:
		var result strings.Builder
		for _, item := range v {
			if str, ok := item.(string); ok {
				result.WriteString(str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				if text, exists := obj["text"]; exists {
					if textStr, ok := text.(string); ok {
						result.WriteString(textStr)
					}
				}
			}
		}
		return result.String()
	default:
		return "[unknown text format]"
	}
}

func (ce *ContextualEmbedder) getCachedEmbedding(text string) []float32 {
	ce.cacheMu.RLock()
	defer ce.cacheMu.RUnlock()

	if embedding, exists := ce.cache[text]; exists {
		return embedding
	}
	return nil
}

func (ce *ContextualEmbedder) setCachedEmbedding(text string, embedding []float32) {
	ce.cacheMu.Lock()
	defer ce.cacheMu.Unlock()

	ce.cache[text] = embedding

	// Периодически сохраняем кэш
	if len(ce.cache)%100 == 0 {
		go ce.saveCache()
	}
}

func (ce *ContextualEmbedder) loadCache() {
	cacheFile := filepath.Join(ce.cacheDir, "embeddings_cache.json")
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return // Кэш файл не существует
	}

	ce.cacheMu.Lock()
	defer ce.cacheMu.Unlock()

	json.Unmarshal(data, &ce.cache)
	log.Printf("[Cache] Загружено %d эмбеддингов из кэша", len(ce.cache))
}

func (ce *ContextualEmbedder) saveCache() {
	if ce.cacheDir == "" {
		return
	}

	cacheFile := filepath.Join(ce.cacheDir, "embeddings_cache.json")

	ce.cacheMu.RLock()
	data, _ := json.Marshal(ce.cache)
	ce.cacheMu.RUnlock()

	os.WriteFile(cacheFile, data, 0644)
}

// Основной класс мигратора
type HistoryMigrator struct {
	storage     *storage.PostgresStorage
	embedder    *ContextualEmbedder
	state       *MigrationState
	stateFile   string
	progressBar *progressbar.ProgressBar
	isPaused    bool
	pauseMu     sync.RWMutex
	maxRetries  int // Максимальное количество попыток для каждого сообщения
}

func NewHistoryMigrator(storage *storage.PostgresStorage, embedder *ContextualEmbedder, stateFile string, maxRetries int) *HistoryMigrator {
	if maxRetries <= 0 {
		maxRetries = 5 // Значение по умолчанию
	}
	
	migrator := &HistoryMigrator{
		storage:    storage,
		embedder:   embedder,
		stateFile:  stateFile,
		state:      &MigrationState{ChatMessagesMap: make(map[int64]int)},
		maxRetries: maxRetries,
	}

	migrator.loadState()
	return migrator
}

func (hm *HistoryMigrator) loadState() {
	data, err := os.ReadFile(hm.stateFile)
	if err != nil {
		log.Printf("[State] Создается новое состояние миграции")
		return
	}

	if err := json.Unmarshal(data, hm.state); err != nil {
		log.Printf("[State] Ошибка загрузки состояния: %v, создается новое", err)
		hm.state = &MigrationState{ChatMessagesMap: make(map[int64]int)}
	} else {
		log.Printf("[State] Загружено состояние: обработано %d/%d сообщений, векторизировано %d",
			hm.state.ProcessedCount, hm.state.TotalMessages, hm.state.VectorizedCount)
	}
}

func (hm *HistoryMigrator) saveState() {
	hm.state.LastUpdateTime = time.Now()
	data, _ := json.MarshalIndent(hm.state, "", "  ")
	os.WriteFile(hm.stateFile, data, 0644)
}

func (hm *HistoryMigrator) Pause() {
	hm.pauseMu.Lock()
	defer hm.pauseMu.Unlock()
	hm.isPaused = true
	log.Printf("[Migration] ПАУЗА активирована")
}

func (hm *HistoryMigrator) Resume() {
	hm.pauseMu.Lock()
	defer hm.pauseMu.Unlock()
	hm.isPaused = false
	log.Printf("[Migration] Миграция ВОЗОБНОВЛЕНА")
}

func (hm *HistoryMigrator) checkPause() {
	for {
		hm.pauseMu.RLock()
		paused := hm.isPaused
		hm.pauseMu.RUnlock()

		if !paused {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (hm *HistoryMigrator) MigrateHistoryDirectory(historyDir string) error {
	// Сканируем все JSON файлы
	files, err := hm.scanJSONFiles(historyDir)
	if err != nil {
		return fmt.Errorf("ошибка сканирования директории: %w", err)
	}

	log.Printf("[Migration] Найдено %d JSON файлов для обработки", len(files))

	// Подсчитываем общее количество сообщений, если еще не подсчитано
	if hm.state.TotalMessages == 0 {
		hm.state.TotalMessages = hm.countTotalMessages(files)
		hm.saveState()
	}

	// Создаем прогресс-бар
	hm.progressBar = progressbar.NewOptions(hm.state.TotalMessages,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan][Migration][reset] Обработка сообщений..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// Устанавливаем текущий прогресс
	hm.progressBar.Set(hm.state.ProcessedCount)

	// Обрабатываем файлы
	for _, file := range files {
		// Проверяем, был ли файл уже обработан
		if hm.isFileProcessed(file) {
			continue
		}

		hm.state.CurrentFile = file
		hm.saveState()

		log.Printf("[Migration] Обработка файла: %s", file)
		if err := hm.processJSONFile(file); err != nil {
			log.Printf("[ERROR] Ошибка обработки файла %s: %v", file, err)
			continue
		}

		hm.state.ProcessedFiles = append(hm.state.ProcessedFiles, file)
		hm.saveState()
	}

	hm.progressBar.Finish()
	log.Printf("[Migration] Миграция завершена! Обработано %d сообщений, векторизировано %d",
		hm.state.ProcessedCount, hm.state.VectorizedCount)

	return nil
}

func (hm *HistoryMigrator) scanJSONFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (hm *HistoryMigrator) countTotalMessages(files []string) int {
	total := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var export TelegramExport
		if err := json.Unmarshal(data, &export); err != nil {
			continue
		}

		for _, msg := range export.Messages {
			if hm.shouldProcessMessage(&msg) {
				total++
			}
		}
	}
	return total
}

func (hm *HistoryMigrator) isFileProcessed(file string) bool {
	for _, processed := range hm.state.ProcessedFiles {
		if processed == file {
			return true
		}
	}
	return false
}

func (hm *HistoryMigrator) shouldProcessMessage(msg *TelegramExportMessage) bool {
	return msg.Type == "message" && len(strings.TrimSpace(hm.embedder.getMessageText(msg))) >= 5
}

func (hm *HistoryMigrator) processJSONFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var export TelegramExport
	if err := json.Unmarshal(data, &export); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	log.Printf("[File] %s: найдено %d сообщений", filepath.Base(filePath), len(export.Messages))

	// Фильтруем сообщения
	var validMessages []TelegramExportMessage
	for _, msg := range export.Messages {
		if hm.shouldProcessMessage(&msg) {
			validMessages = append(validMessages, msg)
		}
	}

	log.Printf("[File] %s: валидных сообщений для обработки: %d", filepath.Base(filePath), len(validMessages))

	// Обрабатываем сообщения
	batchSize := 50
	for i := 0; i < len(validMessages); i += batchSize {
		hm.checkPause() // Проверяем паузу

		end := i + batchSize
		if end > len(validMessages) {
			end = len(validMessages)
		}

		batch := validMessages[i:end]
		if err := hm.processBatch(batch, validMessages, export.ID); err != nil {
			log.Printf("[ERROR] Ошибка обработки пакета %d-%d: %v", i, end, err)
			continue
		}

		// Небольшая пауза между пакетами
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (hm *HistoryMigrator) processBatch(batch []TelegramExportMessage, allMessages []TelegramExportMessage, chatID int64) error {
	for _, msg := range batch {
		hm.checkPause() // Проверяем паузу

		// Пытаемся обработать сообщение с retry
		if err := hm.processMessageWithRetry(&msg, allMessages, chatID); err != nil {
			log.Printf("❌ [CRITICAL] Не удалось обработать сообщение %d после всех попыток: %v", msg.ID, err)
			// Пропускаем это сообщение и продолжаем
			continue
		}

		// Увеличиваем счетчик только после успешной обработки
		hm.state.ProcessedCount++

		// Обновляем прогресс-бар
		if hm.progressBar != nil {
			hm.progressBar.Set(hm.state.ProcessedCount)
		}

		// Сохраняем состояние каждые 10 сообщений
		if hm.state.ProcessedCount%10 == 0 {
			hm.saveState()
		}

		// Детальная статистика каждые 100 сообщений
		if hm.state.ProcessedCount%100 == 0 {
			hm.logDetailedProgress()
		}
	}

	return nil
}

// processMessageWithRetry обрабатывает одно сообщение с retry логикой
func (hm *HistoryMigrator) processMessageWithRetry(msg *TelegramExportMessage, allMessages []TelegramExportMessage, chatID int64) error {
	baseDelay := 2 * time.Second

	var lastErr error

	for attempt := 0; attempt < hm.maxRetries; attempt++ {
		if attempt > 0 {
			// Экспоненциальная задержка: 2s, 4s, 8s, 16s, 32s
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			if delay > 60*time.Second {
				delay = 60 * time.Second
			}
			log.Printf("⏳ [Retry %d/%d] Ожидание %v перед повторной попыткой для сообщения %d", 
				attempt+1, hm.maxRetries, delay, msg.ID)
			time.Sleep(delay)
		}

		// Преобразуем в формат tgbotapi.Message
		tgMsg := hm.convertToTelegramMessage(msg, chatID)

		// Находим индекс сообщения в контексте
		globalIndex := hm.findMessageIndex(msg.ID, allMessages)
		if globalIndex < 0 {
			return fmt.Errorf("message not found in context")
		}

		// Генерируем эмбеддинг
		embedding, context, err := hm.embedder.GenerateEmbeddingWithContext(msg, allMessages, globalIndex)
		if err != nil {
			lastErr = fmt.Errorf("ошибка генерации эмбеддинга: %w", err)
			log.Printf("❌ [Attempt %d/%d] Ошибка генерации эмбеддинга для сообщения %d: %v", 
				attempt+1, hm.maxRetries, msg.ID, err)
			continue
		}

		// Сохраняем сообщение в базу данных
		if err := hm.storage.AddMessage(chatID, tgMsg); err != nil {
			lastErr = fmt.Errorf("ошибка сохранения сообщения в БД: %w", err)
			log.Printf("❌ [Attempt %d/%d] Ошибка сохранения сообщения %d в БД: %v", 
				attempt+1, hm.maxRetries, msg.ID, err)
			continue
		}

		// Сохраняем эмбеддинг в PostgreSQL
		if err := hm.storage.UpdateMessageEmbeddingWithContext(chatID, msg.ID, embedding, context); err != nil {
			lastErr = fmt.Errorf("ошибка сохранения эмбеддинга в БД: %w", err)
			log.Printf("❌ [Attempt %d/%d] Ошибка сохранения эмбеддинга для сообщения %d: %v", 
				attempt+1, hm.maxRetries, msg.ID, err)
			continue
		}

		// Успешно обработано!
		hm.state.VectorizedCount++
		
		if attempt > 0 {
			log.Printf("✅ [Success] Сообщение %d успешно обработано после %d попыток", msg.ID, attempt+1)
		}
		
		return nil
	}

	// Все попытки исчерпаны
	return fmt.Errorf("превышено количество попыток (%d): %w", hm.maxRetries, lastErr)
}

// logDetailedProgress выводит детальную информацию о прогрессе
func (hm *HistoryMigrator) logDetailedProgress() {
	min, day, delay := hm.embedder.rateLimiter.GetStatus()

	// Вычисляем процент выполнения
	percentComplete := float64(hm.state.ProcessedCount) / float64(hm.state.TotalMessages) * 100

	// Вычисляем скорость обработки
	elapsed := time.Since(hm.state.LastUpdateTime)
	if elapsed > 0 && hm.state.ProcessedCount > 0 {
		msgPerSecond := float64(hm.state.ProcessedCount) / elapsed.Seconds()

		// Оценка времени до завершения
		remaining := hm.state.TotalMessages - hm.state.ProcessedCount
		estimatedTimeLeft := time.Duration(float64(remaining)/msgPerSecond) * time.Second

		log.Printf("📊 [Progress] Обработано: %d/%d (%.1f%%) | Векторизировано: %d | Скорость: %.1f сообщ/сек",
			hm.state.ProcessedCount, hm.state.TotalMessages, percentComplete, hm.state.VectorizedCount, msgPerSecond)
		log.Printf("⏱️ [Progress] Осталось времени: ~%v | API лимиты: %d/мин, %d/день | Задержка: %v",
			estimatedTimeLeft.Round(time.Minute), min, day, delay)
	} else {
		log.Printf("📊 [Progress] Обработано: %d/%d (%.1f%%) | Векторизировано: %d | API: %d/мин, %d/день, задержка: %v",
			hm.state.ProcessedCount, hm.state.TotalMessages, percentComplete, hm.state.VectorizedCount, min, day, delay)
	}
}

func (hm *HistoryMigrator) findMessageIndex(messageID int, messages []TelegramExportMessage) int {
	for i, msg := range messages {
		if msg.ID == messageID {
			return i
		}
	}
	return -1
}

func (hm *HistoryMigrator) convertToTelegramMessage(msg *TelegramExportMessage, chatID int64) *tgbotapi.Message {
	// Парсим дату
	var unixTime int64
	if msg.DateUnixtime != "" {
		fmt.Sscanf(msg.DateUnixtime, "%d", &unixTime)
	} else if msg.Date != "" {
		if t, err := time.Parse("2006-01-02T15:04:05", msg.Date); err == nil {
			unixTime = t.Unix()
		}
	}

	// Создаем пользователя
	user := &tgbotapi.User{
		UserName:  msg.From,
		FirstName: msg.From,
	}

	// Парсим FromID если есть
	if msg.FromID != "" {
		fmt.Sscanf(msg.FromID, "user%d", &user.ID)
	}

	// Создаем чат
	chat := &tgbotapi.Chat{
		ID: chatID,
	}

	// Создаем сообщение
	tgMsg := &tgbotapi.Message{
		MessageID:      msg.ID,
		From:           user,
		Date:           int(unixTime),
		Chat:           chat,
		Text:           hm.embedder.getMessageText(msg),
		ReplyToMessage: nil, // Можно доработать для reply chains
	}

	// Добавляем reply если есть
	if msg.ReplyToMessageID != 0 {
		tgMsg.ReplyToMessage = &tgbotapi.Message{MessageID: msg.ReplyToMessageID}
	}

	return tgMsg
}

func (hm *HistoryMigrator) updateMessageEmbedding(chatID int64, messageID int, embedding []float32, context string) error {
	// Здесь нужно будет добавить метод в PostgresStorage для обновления эмбеддинга
	// Пока заглушка - в следующем файле создадим расширение storage
	log.Printf("[DEBUG] Нужно сохранить эмбеддинг для сообщения %d (размер: %d)", messageID, len(embedding))
	return nil
}

// Утилиты
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	log.Printf("[Migration] 🚀 Запуск миграции истории чата...")

	// Определяем рабочую директорию скрипта
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("❌ Ошибка определения пути скрипта: %v", err)
	}
	workDir := filepath.Dir(execPath)

	// Если запускаем через go run, используем текущую директорию
	if strings.Contains(execPath, "tmp") || strings.Contains(execPath, "go-build") {
		workDir, _ = os.Getwd()
	}

	log.Printf("[Migration] 📁 Рабочая директория: %s", workDir)

	// Загружаем .env файл из рабочей директории
	envPath := filepath.Join(workDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := utils.LoadEnvFile(envPath); err != nil {
			log.Printf("⚠️ Предупреждение: не удалось загрузить %s: %v", envPath, err)
		} else {
			log.Printf("[Migration] ✅ Загружен конфигурационный файл: %s", envPath)
		}
	} else {
		log.Fatalf("❌ Не найден обязательный файл .env в директории: %s", workDir)
	}

	// Загружаем только необходимые переменные окружения для миграции
	dbHost := os.Getenv("POSTGRESQL_HOST")
	dbPort := os.Getenv("POSTGRESQL_PORT")
	dbUser := os.Getenv("POSTGRESQL_USER")
	dbPassword := os.Getenv("POSTGRESQL_PASSWORD")
	dbName := os.Getenv("POSTGRESQL_DBNAME")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// Настройки для векторизации с значениями по умолчанию
	embeddingModel := utils.GetEnvOrDefault("GEMINI_EMBEDDING_MODEL_NAME", "embedding-001")

	// Пути относительно рабочей директории
	historyDir := filepath.Join(workDir, "data")
	cacheDir := filepath.Join(workDir, "cache", "embeddings")
	stateFile := filepath.Join(workDir, "migration_state.json")
	logFile := filepath.Join(workDir, "migration.log")

	// Проверяем обязательные переменные
	requiredVars := map[string]string{
		"POSTGRESQL_HOST":     dbHost,
		"POSTGRESQL_USER":     dbUser,
		"POSTGRESQL_PASSWORD": dbPassword,
		"POSTGRESQL_DBNAME":   dbName,
		"GEMINI_API_KEY":      geminiAPIKey,
	}

	for varName, value := range requiredVars {
		if value == "" {
			log.Fatalf("❌ Обязательная переменная %s не задана в .env", varName)
		}
	}

	if dbPort == "" {
		dbPort = "5432"
	}

	// Настраиваем логирование в файл
	setupFileLogging(logFile)

	log.Printf("[Migration] 🔗 Подключение к PostgreSQL: %s:%s (База: %s)", dbHost, dbPort, dbName)
	log.Printf("[Migration] 📊 Модель эмбеддингов: %s", embeddingModel)
	log.Printf("[Migration] 📂 Директория с историей: %s", historyDir)
	log.Printf("[Migration] 💾 Файл состояния: %s", stateFile)
	log.Printf("[Migration] 🗃️ Кэш эмбеддингов: %s", cacheDir)

	// Проверяем существование директории с историей
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		log.Fatalf("❌ Директория с историей не найдена: %s\nСоздайте директорию и поместите в неё JSON файлы экспорта Telegram", historyDir)
	}

	// Инициализируем PostgreSQL storage
	storage, err := storage.NewPostgresStorage(
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		1000, // contextWindow
		true, // debug
	)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к PostgreSQL: %v", err)
	}
	defer storage.Close()
	log.Printf("[Migration] ✅ Успешное подключение к базе данных")

	// Создаем минимальную конфигурацию для Gemini
	cfg := &config.Config{
		GeminiAPIKey:             geminiAPIKey,
		GeminiModelName:          "gemini-2.5-flash",
		GeminiEmbeddingModelName: embeddingModel,
		Debug:                    true,
	}

	// Инициализируем Gemini клиент для эмбеддингов
	llmClient, err := gemini.New(cfg, cfg.GeminiModelName, cfg.GeminiEmbeddingModelName, cfg.Debug)
	if err != nil {
		log.Fatalf("❌ Ошибка инициализации Gemini клиента: %v", err)
	}
	defer llmClient.Close()
	log.Printf("[Migration] ✅ Gemini клиент инициализирован")

	// Парсим конфигурацию для rate limiter из .env
	reqPerMin := utils.ParseIntOrDefault(utils.GetEnvOrDefault("EMBEDDING_REQUESTS_PER_MINUTE", "240"), 240)
	reqPerDay := utils.ParseIntOrDefault(utils.GetEnvOrDefault("EMBEDDING_REQUESTS_PER_DAY", "24000"), 24000)
	delayStr := utils.GetEnvOrDefault("EMBEDDING_REQUEST_DELAY", "300ms")
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		delay = 300 * time.Millisecond
		log.Printf("⚠️ Неверный формат EMBEDDING_REQUEST_DELAY, используется %v", delay)
	}

	log.Printf("[Migration] ⚙️ Rate Limiter: %d req/мин, %d req/день, задержка %v", reqPerMin, reqPerDay, delay)

	// Создаем rate limiter
	rateLimiter := NewRateLimiter(reqPerMin, reqPerDay, delay)

	// Создаем embedder
	embedder := NewContextualEmbedder(llmClient, rateLimiter, cacheDir)

	// Получаем maxRetries из .env
	maxRetries := utils.ParseIntOrDefault(utils.GetEnvOrDefault("EMBEDDING_MAX_RETRIES", "5"), 5)
	log.Printf("[Migration] 🔄 Максимальное количество повторных попыток: %d", maxRetries)

	// Создаем мигратор
	migrator := NewHistoryMigrator(storage, embedder, stateFile, maxRetries)

	// Настраиваем обработку сигналов для корректной остановки
	setupSignalHandling(migrator)

	// Запускаем миграцию
	log.Printf("[Migration] 🎯 Начинаем миграцию из директории: %s", historyDir)
	if err := migrator.MigrateHistoryDirectory(historyDir); err != nil {
		log.Fatalf("❌ Ошибка миграции: %v", err)
	}

	log.Printf("[Migration] 🎉 Миграция успешно завершена!")
}

// setupFileLogging настраивает логирование в файл
func setupFileLogging(logFile string) {
	// Создаем директорию для логов если не существует
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		log.Printf("⚠️ Не удалось создать директорию для логов: %v", err)
		return
	}

	// Открываем файл для записи логов
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("⚠️ Не удалось открыть файл логов %s: %v", logFile, err)
		return
	}

	// Настраиваем мультиплексированный вывод (консоль + файл)
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	log.Printf("[Migration] 📝 Логирование настроено: консоль + %s", logFile)
}

// setupSignalHandling настраивает обработку сигналов для корректной остановки
func setupSignalHandling(migrator *HistoryMigrator) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("[Migration] 🛑 Получен сигнал %v, корректное завершение...", sig)

		// Приостанавливаем миграцию
		migrator.Pause()

		// Сохраняем состояние
		migrator.saveState()

		log.Printf("[Migration] 💾 Состояние сохранено, выход...")
		os.Exit(0)
	}()

	log.Printf("[Migration] 🔧 Обработчик сигналов настроен (Ctrl+C для остановки)")
}
