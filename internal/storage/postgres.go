package storage

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

// PostgresStorage представляет хранилище на основе PostgreSQL
type PostgresStorage struct {
	db            *sql.DB
	contextWindow int
	debug         bool
}

// NewPostgresStorage создает новое PostgreSQL хранилище
func NewPostgresStorage(host, port, user, password, dbname string, contextWindow int, debug bool) (*PostgresStorage, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	storage := &PostgresStorage{
		db:            db,
		contextWindow: contextWindow,
		debug:         debug,
	}

	// Создаем таблицы если не существуют
	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return storage, nil
}

// createTables создает необходимые таблицы
func (ps *PostgresStorage) createTables() error {
	// Создаем расширение для векторов если не существует
	_, err := ps.db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		log.Printf("Warning: could not create vector extension: %v", err)
	}

	// Создаем таблицу сообщений
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chat_id BIGINT NOT NULL,
		message_id INTEGER NOT NULL,
		user_id BIGINT,
		username TEXT,
		first_name TEXT,
		last_name TEXT,
		text TEXT,
		date_sent TIMESTAMP WITH TIME ZONE,
		reply_to_message_id INTEGER,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		embedding vector(768),
		context_text TEXT,
		UNIQUE(chat_id, message_id)
	)`

	_, err = ps.db.Exec(createMessagesTable)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// Создаем индексы
	createIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_date_sent ON messages(date_sent)",
		"CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id)",
	}

	for _, indexSQL := range createIndexes {
		_, err = ps.db.Exec(indexSQL)
		if err != nil {
			log.Printf("Warning: could not create index: %v", err)
		}
	}

	return nil
}

// AddMessage добавляет сообщение в базу данных
func (ps *PostgresStorage) AddMessage(chatID int64, message *tgbotapi.Message) error {
	query := `
		INSERT INTO messages (chat_id, message_id, user_id, username, first_name, last_name, text, date_sent, reply_to_message_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (chat_id, message_id) DO NOTHING
	`

	var userID sql.NullInt64
	var username, firstName, lastName sql.NullString
	var replyToMessageID sql.NullInt32

	if message.From != nil {
		userID = sql.NullInt64{Int64: message.From.ID, Valid: true}
		username = sql.NullString{String: message.From.UserName, Valid: true}
		firstName = sql.NullString{String: message.From.FirstName, Valid: true}
		lastName = sql.NullString{String: message.From.LastName, Valid: true}
	}

	if message.ReplyToMessage != nil {
		replyToMessageID = sql.NullInt32{Int32: int32(message.ReplyToMessage.MessageID), Valid: true}
	}

	_, err := ps.db.Exec(query, chatID, message.MessageID, userID, username, firstName, lastName,
		message.Text, message.Time(), replyToMessageID)

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	if ps.debug {
		log.Printf("Added message %d to chat %d", message.MessageID, chatID)
	}

	return nil
}

// UpdateMessageEmbeddingWithContext обновляет эмбеддинг сообщения с контекстом
func (ps *PostgresStorage) UpdateMessageEmbeddingWithContext(chatID int64, messageID int, embedding []float32, context string) error {
	query := `
		UPDATE messages 
		SET embedding = $1, context_text = $2
		WHERE chat_id = $3 AND message_id = $4
	`

	// Конвертируем []float32 в строку для pgvector
	embeddingStr := fmt.Sprintf("[%v]", embedding)

	_, err := ps.db.Exec(query, embeddingStr, context, chatID, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message embedding: %w", err)
	}

	if ps.debug {
		log.Printf("Updated embedding for message %d in chat %d", messageID, chatID)
	}

	return nil
}

// Close закрывает соединение с базой данных
func (ps *PostgresStorage) Close() error {
	if ps.db != nil {
		return ps.db.Close()
	}
	return nil
}
