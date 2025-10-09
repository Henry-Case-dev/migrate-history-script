# 🚀 Migrate History Script

Скрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.

## ⚡ Быстрый старт

### 1. Клонирование и установка

```bash
git clone https://github.com/Henry-Case-dev/migrate-history-script.git
cd migrate-history-script
go mod download
```

### 2. Настройка окружения

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

# Настройки векторизации (опционально)
EMBEDDING_REQUESTS_PER_MINUTE=240
EMBEDDING_REQUESTS_PER_DAY=24000
EMBEDDING_REQUEST_DELAY=0.3s
EMBEDDING_BATCH_SIZE=100
EMBEDDING_CACHE_ENABLED=true

# Retry настройки
EMBEDDING_MAX_RETRIES=5  # 0 = бесконечный режим с 24ч таймаутом
```

### 3. Подготовка данных

Поместите JSON файлы экспорта Telegram в папку `data/`:

```
data/
├── chat_export_1.json
├── chat_export_2.json
└── ...
```

### 4. Запуск

```bash
# Проверка настроек
go run cmd/test_setup/main.go

# Запуск миграции
go run main.go

# Или скомпилированной версией
go build -o migrate_history.exe
./migrate_history.exe
```

## 📋 Основные возможности

- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры
- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  
- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных
- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени
- ✅ **Rate limiting** - защита от превышения квот API
- ✅ **Кэширование эмбеддингов** - экономия API вызовов
- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

## 📋 Требования

- **Go 1.21+** 
- **PostgreSQL** с установленным pgvector
- **Gemini API ключ** (tier 1 рекомендуется)
- **~500MB** свободного места в БД
- **~$10-15** для векторизации большого архива

## 🔧 Продвинутые настройки

### Переменные окружения

```env
# Производительность
EMBEDDING_BATCH_SIZE=100
EMBEDDING_BATCH_DELAY=60s
EMBEDDING_REQUESTS_PER_MINUTE=240
EMBEDDING_REQUEST_DELAY=0.3s

# Retry механизм
EMBEDDING_MAX_RETRIES=5        # Количество попыток при ошибках
                               # 0 = бесконечный режим с 24ч таймаутом

# Кэширование
EMBEDDING_CACHE_ENABLED=true
EMBEDDING_CACHE_DIR=./cache/embeddings

# Возобновление миграции
MIGRATION_RESUME_ENABLED=true
MIGRATION_STATE_FILE=./migration_state.json
```

### Retry механизм

Скрипт поддерживает два режима обработки ошибок:

#### Ограниченный режим (EMBEDDING_MAX_RETRIES > 0)
- Делает указанное количество попыток при ошибках
- Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (максимум 60с)
- При исчерпании попыток возвращает ошибку

**Пример**: `EMBEDDING_MAX_RETRIES=5` - 5 попыток, затем ошибка

#### Бесконечный режим (EMBEDDING_MAX_RETRIES = 0)
- **Никогда не пропускает сообщения** - гарантия полноты данных
- Повторяет попытки обработки до успеха или таймаута
- Экспоненциальная задержка до **60 минут** между попытками
- После достижения 60 минут - фиксированные интервалы в 60 минут
- **Максимальная длительность**: 24 часа на одно сообщение
- При таймауте - **завершение программы** с сохранением состояния

**Рекомендуется для production**: гарантирует что все сообщения будут обработаны

**Когда использовать:**
- `EMBEDDING_MAX_RETRIES=5` - для тестирования и быстрой миграции
- `EMBEDDING_MAX_RETRIES=0` - для production, когда важна полнота данных

### Контекстная векторизация

Скрипт использует адаптивные контекстные окна:

- **Базовое окно**: 5+5 сообщений для обычных сообщений
- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  
- **Максимальное окно**: 25+25 для важных дискуссий
- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

## 🚦 Управление процессом

### Запуск и остановка

```bash
# Graceful остановка
Ctrl+C

# Возобновление (автоматически продолжит с места остановки)
./migrate_history.exe

# Запуск в фоне
nohup ./migrate_history.exe > migration.log 2>&1 &
```

### Мониторинг

```bash
# Логи в реальном времени
tail -f migration.log

# Статус миграции
cat migration_state.json | jq '.'

# API квоты
grep "Rate limit\|API:" migration.log
```

## 🛡️ Безопасность при переносе

При переносе скрипта на другую машину **уже обработанные данные не будут загружаться повторно** благодаря:

1. **Отслеживание состояния** - файл `migration_state.json` хранит информацию о обработанных файлах
2. **Database constraints** - PostgreSQL UNIQUE ограничения на `(chat_id, message_id)`
3. **ON CONFLICT стратегия** - `INSERT ... ON CONFLICT DO NOTHING` предотвращает дублирование

**Для переноса:**

1. Скопируйте `migration_state.json` 
2. Убедитесь что PostgreSQL содержит уже импортированные данные
3. Запустите скрипт - он автоматически пропустит обработанные файлы

## 📊 Структура базы данных

Скрипт автоматически создает таблицу `messages`:

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    message_id INTEGER NOT NULL,
    user_id BIGINT,
    username TEXT,
    text TEXT,
    date_sent TIMESTAMP WITH TIME ZONE,
    embedding vector(768),
    context_text TEXT,
    UNIQUE(chat_id, message_id)
);
```

### Проверка результата

```sql
-- Общая статистика
SELECT 
    COUNT(*) as total_messages,
    COUNT(embedding) as vectorized_messages,
    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as vectorization_percent
FROM messages;

-- Топ чаты по количеству сообщений
SELECT chat_id, COUNT(*) as messages
FROM messages 
GROUP BY chat_id 
ORDER BY messages DESC 
LIMIT 10;
```

## ⚡ Производительность

### Rate Limiting

- **Tier 1 Gemini**: 240 запросов/мин, 24,000/день
- **Адаптивное замедление** при достижении лимитов
- **Экспоненциальная задержка** при ошибках

### Оптимизация

- **Пакетная обработка**: 50-100 сообщений за раз
- **Локальный кэш** эмбеддингов для ускорения
- **Дедупликация** одинаковых контекстов

## 🔧 Решение проблем

### Ошибки подключения к PostgreSQL

```bash
# Проверка подключения
psql -h your_host -p 5432 -U your_user -d your_database
```

### Превышение лимитов API

- Подождите сброса лимитов (обновляется каждую минуту/день)
- Уменьшите `EMBEDDING_REQUESTS_PER_MINUTE` в .env
- Используйте резервный API ключ

### 24-часовой таймаут (EMBEDDING_MAX_RETRIES=0)

Если программа завершилась с ошибкой "Не удалось обработать сообщение за 24 часа":

1. **Проверьте Gemini API** - возможно проблемы на стороне сервиса
2. **Проверьте API ключ** - убедитесь что квоты не исчерпаны
3. **Проверьте подключение** - стабильность интернет соединения
4. **После устранения проблемы** просто перезапустите скрипт - он продолжит с места остановки

### Проблемы с памятью

- Уменьшите `EMBEDDING_BATCH_SIZE`
- Увеличьте `EMBEDDING_BATCH_DELAY`  
- Перезапустите процесс

### Ошибки парсинга JSON

- Проверьте кодировку файлов (должна быть UTF-8)
- Убедитесь в валидности JSON структуры
- Исключите поврежденные файлы

## 📁 Структура проекта

```
migrate-history-script/
├── main.go                    # Основной файл программы
├── internal/                  # Внутренние пакеты
│   ├── config/               # Конфигурация
│   ├── storage/              # Работа с PostgreSQL
│   ├── gemini/               # Клиент Gemini API
│   ├── llm/                  # Интерфейсы LLM
│   └── utils/                # Утилиты
├── cmd/
│   └── test_setup/           # Скрипт проверки настроек
├── data/                     # JSON файлы экспорта Telegram
├── cache/                    # Кэш эмбеддингов
├── migration_state.json     # Состояние миграции (автосоздается)
└── .env                      # Конфигурация (создать самостоятельно)
```

## 🎯 После миграции

1. **Переключение бота** на PostgreSQL storage
2. **Тестирование векторного поиска** 
3. **Настройка индексов** для оптимизации
4. **Мониторинг производительности**

## 📈 Векторный поиск

После миграции вы сможете использовать семантический поиск:

```sql
-- Поиск похожих сообщений
SELECT 
    text,
    embedding <=> '[vector_here]'::vector as distance
FROM messages 
WHERE embedding IS NOT NULL
ORDER BY distance ASC 
LIMIT 10;

-- Создание индекса для быстрого поиска
CREATE INDEX ON messages USING ivfflat (embedding vector_cosine_ops);
```

## 📞 Поддержка

При возникновении проблем:

1. Запустите `go run cmd/test_setup/main.go` для диагностики
2. Проверьте логи в файле `migration.log`
3. Убедитесь в корректности настроек в `.env` файле

---

**Автор**: Henry-Case-dev  
**Репозиторий**: https://github.com/Henry-Case-dev/migrate-history-script  
**Версия**: 1.0.0  
**Последнее обновление**: Сентябрь 2025