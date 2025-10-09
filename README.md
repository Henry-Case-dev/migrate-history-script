# � Telegram History Migration Tool

> Автоматическая миграция истории Telegram чатов в PostgreSQL с AI-векторизацией через Gemini API

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## ⚡ Быстрый старт

### 1. Установка Go (если еще не установлен)

**Windows:**
```bash
# Скачайте установщик: https://go.dev/dl/
# Или через Chocolatey:
choco install golang
```

**Linux/Mac:**
```bash
# Linux
sudo apt-get install golang-go

# Mac
brew install go
```

Проверка: `go version` (должно показать Go 1.21+)

### 2. Клонирование и установка зависимостей

```bash
git clone https://github.com/Henry-Case-dev/migrate-history-script.git
cd migrate-history-script
go mod download
```

### 3. Настройка окружения

Создайте файл `.env` в корне проекта:

```env
# PostgreSQL
POSTGRESQL_HOST=your_host
POSTGRESQL_PORT=5432
POSTGRESQL_USER=your_user
POSTGRESQL_PASSWORD=your_password
POSTGRESQL_DBNAME=your_database

# Gemini API
GEMINI_API_KEY=your_gemini_api_key
GEMINI_EMBEDDING_MODEL_NAME=embedding-001

# Retry настройки
EMBEDDING_MAX_RETRIES=0  # 0 = бесконечный режим с 24ч таймаутом

# Rate Limiting (рекомендуемые значения для Tier 1)
EMBEDDING_REQUESTS_PER_MINUTE=240
EMBEDDING_REQUESTS_PER_DAY=24000
EMBEDDING_REQUEST_DELAY=0.3s
```

### 4. Подготовка данных

Поместите JSON файлы экспорта Telegram в папку `data/`:

```
data/
├── chat_export_1.json
├── chat_export_2.json
└── ...
```

### 5. Запуск

```bash
# Проверка настроек (опционально)
go run cmd/test_setup/main.go

# Запуск миграции
go run main.go

# Или скомпилируйте и запустите
go build -o migrate_history.exe
./migrate_history.exe
```

### 6. Остановка

```bash
# Graceful остановка: Ctrl+C
# Программа сохранит состояние

# При следующем запуске автоматически продолжит:
./migrate_history.exe
```

---

## � Основные возможности

- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры
- ✅ **AI-векторизация** - умное создание эмбеддингов через Gemini API
- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных
- ✅ **Бесконечный retry режим** - гарантия обработки всех сообщений
- ✅ **Защита от дубликатов** - автоматический пропуск уже обработанных
- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени
- ✅ **Rate limiting** - защита от превышения квот API
- ✅ **Кэширование эмбеддингов** - экономия API вызовов

---

## � Требования

| Компонент | Версия | Примечание |
|-----------|--------|------------|
| **Go** | 1.21+ | Язык программирования |
| **PostgreSQL** | 13+ с pgvector | База данных с векторным расширением |
| **Gemini API** | Tier 1+ | Для генерации эмбеддингов (240 req/min) |
| **Дисковое пространство** | ~500MB | Для хранения данных |

### Установка pgvector

```sql
-- В PostgreSQL выполните:
CREATE EXTENSION IF NOT EXISTS vector;
```

---

## � Режимы работы

### �️ Бесконечный режим (РЕКОМЕНДУЕТСЯ)

```env
EMBEDDING_MAX_RETRIES=0
```

**Поведение:**
- ✅ Никогда не пропускает сообщения
- ✅ Повторяет попытки до успеха или таймаута (24 часа)
- ✅ Экспоненциальная задержка до 60 минут между попытками
- ✅ При таймауте - завершение с сохранением состояния

**Когда использовать:** production, когда важна полнота данных

### � Ограниченный режим

```env
EMBEDDING_MAX_RETRIES=5
```

**Поведение:**
- ✅ Делает 5 попыток при ошибках
- ✅ Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (max 60с)
- ⚠️ При исчерпании попыток - возвращает ошибку

**Когда использовать:** тестирование, быстрая миграция

---

## � Управление процессом

### Запуск и остановка

```bash
# Запуск
./migrate_history.exe

# Graceful остановка
Ctrl+C

# Возобновление
./migrate_history.exe

# Запуск в фоне (Linux/Mac)
nohup ./migrate_history.exe > migration.log 2>&1 &

# Windows (PowerShell)
Start-Process -NoNewWindow -FilePath ".\migrate_history.exe" -RedirectStandardOutput "migration.log"
```

### Мониторинг

```bash
# Логи в реальном времени
tail -f migration.log

# Статус миграции
cat migration_state.json | jq '.'

# Проверка квот API
grep "Rate limit\|API:" migration.log
```

### Проверка результатов в БД

```sql
-- Общая статистика
SELECT 
    COUNT(*) as total_messages,
    COUNT(embedding) as vectorized_messages,
    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as percent
FROM messages;

-- Топ-10 чатов
SELECT chat_id, COUNT(*) as messages
FROM messages 
GROUP BY chat_id 
ORDER BY messages DESC 
LIMIT 10;
```

---

## �️ Защита от дубликатов

При переносе скрипта на другую машину **данные не дублируются**:

1. ✅ `migration_state.json` - отслеживает обработанные файлы и индексы
2. ✅ PostgreSQL UNIQUE constraint на `(chat_id, message_id)`
3. ✅ `INSERT ... ON CONFLICT DO NOTHING` стратегия

**Для переноса:**

```bash
# 1. Скопируйте файл состояния
scp migration_state.json user@server:/path/

# 2. Убедитесь что PostgreSQL доступна
psql -h your_host -U your_user -d your_database -c "SELECT COUNT(*) FROM messages;"

# 3. Запустите - программа продолжит с места остановки
./migrate_history.exe
```

---

## � Решение проблем

### ❌ Ошибка подключения к PostgreSQL

```bash
# Проверьте подключение
psql -h your_host -p 5432 -U your_user -d your_database

# Проверьте расширение pgvector
psql -h your_host -U your_user -d your_database -c "SELECT * FROM pg_extension WHERE extname='vector';"
```

**Решение:**
- Проверьте `.env` настройки (host, port, user, password)
- Убедитесь что PostgreSQL запущен
- Установите pgvector: `CREATE EXTENSION vector;`

### ⏱️ Таймаут 24 часа (EMBEDDING_MAX_RETRIES=0)

**Причины:**
- Проблемы с Gemini API (перегрузка сервиса)
- Исчерпание квот API
- Нестабильное интернет-соединение

**Решение:**
1. Проверьте статус Gemini API
2. Проверьте квоты: https://aistudio.google.com/app/apikey
3. Проверьте интернет-соединение
4. После устранения - просто перезапустите скрипт

### � Превышение лимитов API

```env
# Решение 1: Уменьшите лимиты в .env
EMBEDDING_REQUESTS_PER_MINUTE=120  # Было 240
EMBEDDING_REQUEST_DELAY=0.6s       # Было 0.3s

# Решение 2: Подождите сброса (автоматически)
# Tier 1: обновляется каждую минуту/день

# Решение 3: Используйте резервный ключ
GEMINI_API_KEY=your_backup_key
```

---

## � Структура базы данных

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    message_id INTEGER NOT NULL,
    user_id BIGINT,
    username TEXT,
    text TEXT,
    date_sent TIMESTAMP WITH TIME ZONE,
    embedding vector(768),        -- Gemini embedding-001
    context_text TEXT,
    UNIQUE(chat_id, message_id)
);

-- Индекс для векторного поиска
CREATE INDEX messages_embedding_idx 
ON messages 
USING ivfflat (embedding vector_cosine_ops);
```

---

## � Векторный поиск (после миграции)

```sql
-- Поиск похожих сообщений (семантический поиск)
SELECT 
    text,
    embedding <=> '[your_vector]'::vector as distance
FROM messages 
WHERE embedding IS NOT NULL
ORDER BY distance ASC 
LIMIT 10;

-- Поиск в конкретном чате
SELECT text, embedding <=> '[vector]'::vector as dist
FROM messages 
WHERE chat_id = 123456 AND embedding IS NOT NULL
ORDER BY dist ASC 
LIMIT 5;
```

---

## � Структура проекта

```
migrate-history-script/
├── main.go                    # Основной файл (миграция, retry логика)
├── internal/
│   ├── config/                # Загрузка .env конфигурации
│   ├── storage/               # PostgreSQL клиент
│   ├── gemini/                # Gemini API клиент
│   ├── llm/                   # LLM интерфейсы
│   └── utils/                 # Утилиты
├── cmd/
│   └── test_setup/            # Проверка настроек БД/API
├── data/                      # JSON экспорты Telegram
├── cache/                     # Кэш эмбеддингов
├── migration_state.json       # Состояние миграции
└── .env                       # Конфигурация (создайте вручную)
```

---

## ⚡ Производительность

### Rate Limiting

| Tier | Requests/min | Requests/day | Рекомендации |
|------|--------------|--------------|--------------|
| **Free** | 15 | 1,500 | Не подходит для больших архивов |
| **Tier 1** | 240 | 24,000 | ✅ Рекомендуется |
| **Tier 2** | 1,000+ | 100,000+ | Для очень больших архивов |

### Оптимизация

- ⚡ Пакетная обработка (50 сообщений в батче)
- � Локальный кэш эмбеддингов
- � Дедупликация одинаковых контекстов
- � Адаптивное замедление при лимитах
- � Возобновление с точного checkpoint

---

## � Поддержка

**Диагностика проблем:**

```bash
# 1. Запустите проверку настроек
go run cmd/test_setup/main.go

# 2. Проверьте логи
tail -n 100 migration.log

# 3. Проверьте состояние
cat migration_state.json | jq '.'
```

**Полезные ресурсы:**
- [Gemini API Documentation](https://ai.google.dev/docs)
- [PostgreSQL pgvector](https://github.com/pgvector/pgvector)
- [Issues](https://github.com/Henry-Case-dev/migrate-history-script/issues)

---

## � Лицензия

MIT License - см. [LICENSE](LICENSE)

---

<div align="center">

**Автор:** [Henry-Case-dev](https://github.com/Henry-Case-dev)  
**Репозиторий:** [migrate-history-script](https://github.com/Henry-Case-dev/migrate-history-script)  

⭐ Поставьте звезду, если проект был полезен!

</div>
