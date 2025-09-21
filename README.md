# Migrate History Script# 🚀 Миграция истории чата в PostgreSQL с векторизацией



Скрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.## Описание



## Быстрый стартЭтот инструмент выполняет миграцию истории чата из JSON файлов в PostgreSQL с автоматической векторизацией для поиска по семантическому сходству.



### 1. Клонирование проекта## Основные возможности



```bash- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры

git clone <repository-url>- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  

cd migrate-history-script- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных

```- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени

- ✅ **Rate limiting** - защита от превышения квот API

### 2. Установка зависимостей- ✅ **Кэширование эмбеддингов** - экономия API вызовов

- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

```bash

go mod download## Требования

```

- PostgreSQL с установленным pgvector

### 3. Настройка конфигурации- Действующий API ключ Gemini (tier 1)

- Go 1.21+

Создайте файл `.env` в корне проекта:- ~500MB свободного места в БД

- ~$10-15 для векторизации

```env

# PostgreSQL## Настройка

POSTGRESQL_HOST=your_host

POSTGRESQL_PORT=5432### 1. Переменные окружения (.env)

POSTGRESQL_USER=your_user

POSTGRESQL_PASSWORD=your_password```env

POSTGRESQL_DBNAME=your_database# PostgreSQL (уже настроено)

POSTGRESQL_HOST=46.19.69.26

# Gemini APIPOSTGRESQL_PORT=5432

GEMINI_API_KEY=your_gemini_api_keyPOSTGRESQL_USER=gen_user

GEMINI_EMBEDDING_MODEL_NAME=embedding-001POSTGRESQL_PASSWORD=poshelnahuy0880

POSTGRESQL_DBNAME=default_db

# Настройки производительности (опционально)

EMBEDDING_REQUESTS_PER_MINUTE=240# Настройки векторизации

EMBEDDING_REQUESTS_PER_DAY=24000EMBEDDING_REQUESTS_PER_MINUTE=240

EMBEDDING_REQUEST_DELAY=300msEMBEDDING_REQUESTS_PER_DAY=24000

```EMBEDDING_BATCH_SIZE=100

EMBEDDING_REQUEST_DELAY=0.3s

### 4. Подготовка данныхEMBEDDING_BATCH_DELAY=60s

EMBEDDING_ADAPTIVE_THROTTLING=true

Поместите JSON файлы экспорта Telegram в папку `data/`:EMBEDDING_SAFETY_MARGIN=0.8

EMBEDDING_MAX_RETRIES=3

```EMBEDDING_CACHE_ENABLED=true

data/EMBEDDING_CACHE_DIR=./cache/embeddings

├── chat_export_1.jsonMIGRATION_RESUME_ENABLED=true

├── chat_export_2.jsonMIGRATION_STATE_FILE=./migration_state.json

└── ...

```# Gemini API (проверьте свой ключ)

GEMINI_API_KEY=your_api_key_here

### 5. ЗапускGEMINI_EMBEDDING_MODEL_NAME=embedding-001

```

#### Проверка настроек:

```bash### 2. Структура директорий

go run cmd/test_setup/main.go

``````

rofloslav/

#### Запуск миграции:├── cmd/migrate_history/

```bash│   ├── migrate_history.exe     # Скомпилированный инструмент

go run main.go│   ├── run_migration.sh       # Скрипт запуска

```│   └── main.go               # Исходный код

└── История чата/             # JSON файлы истории

или скомпилированной версией:    ├── chat_history_1.json

```bash    ├── chat_history_2.json

go build -o migrate_history.exe    └── ...

./migrate_history.exe```

```

## Использование

## Возможности

### Быстрый старт

- ✅ **Импорт JSON → PostgreSQL** - перенос истории с сохранением структуры

- ✅ **Контекстная векторизация** - умные эмбеддинги с учетом соседних сообщений  ```bash

- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данныхcd cmd/migrate_history

- ✅ **Защита от дублирования** - безопасный повторный запуск./run_migration.sh

- ✅ **Rate limiting** - защита от превышения квот API```

- ✅ **Кэширование эмбеддингов** - экономия API вызовов

# В фоне:

## Структура проекта```bash

nohup ./migrate_history.exe > migration.log 2>&1 &

``````

migrate-history-script/

├── main.go                    # Основной файл программы### Ручной запуск

├── internal/                  # Внутренние пакеты

│   ├── config/               # Конфигурация```bash

│   ├── storage/              # Работа с PostgreSQLcd cmd/migrate_history

│   ├── gemini/               # Клиент Gemini API./migrate_history.exe

│   ├── llm/                  # Интерфейсы LLM```

│   └── utils/                # Утилиты

├── cmd/### Возобновление прерванной миграции

│   └── test_setup/           # Скрипт проверки настроек

├── data/                     # JSON файлы экспорта TelegramПросто запустите инструмент снова - он автоматически продолжит с места остановки используя `migration_state.json`.

├── cache/                    # Кэш эмбеддингов

├── migration_state.json     # Состояние миграции (автосоздается)## Стратегия векторизации

└── .env                      # Конфигурация (создать самостоятельно)

```### Контекстные окна



## Управление процессом- **Базовое окно**: 5+5 сообщений для обычных сообщений

- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  

- **Остановка**: `Ctrl+C` (состояние сохраняется автоматически)- **Максимальное окно**: 25+25 для важных дискуссий

- **Возобновление**: Перезапустите программу - автоматически продолжит с места остановки- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

- **Мониторинг**: Логи выводятся в консоль и файл `migration.log`

### Адаптивная логика

## Безопасность при переносе

```go

При переносе скрипта на другую машину:// Определение размера окна

if isImportantMessage(msg) {        // Длинные сообщения, вопросы

1. Скопируйте файл `migration_state.json`     return 25                       // Максимальный контекст

2. Убедитесь что PostgreSQL содержит уже импортированные данные} else if hasReplyChain(msg) {      // Ответы в цепочке

3. Скрипт автоматически пропустит уже обработанные файлы и сообщения    return 15                       // Средний контекст  

} else if containsQuestions(msg) {   // Вопросы

Механизм защиты:    return 10                       // Увеличенный контекст

- Отслеживание обработанных файлов в `migration_state.json`} else {

- Уникальные constraint'ы в PostgreSQL: `UNIQUE(chat_id, message_id)`    return 5                        // Базовый контекст

- `INSERT ... ON CONFLICT DO NOTHING` стратегия}

```

## База данных

## Мониторинг

Скрипт автоматически создает таблицу `messages`:

### Логи в реальном времени

```sql

CREATE TABLE messages (```bash

    id SERIAL PRIMARY KEY,tail -f migration.log

    chat_id BIGINT NOT NULL,```

    message_id INTEGER NOT NULL,

    user_id BIGINT,### Статус миграции

    username TEXT,

    text TEXT,```bash

    date_sent TIMESTAMP WITH TIME ZONE,cat migration_state.json | jq '.'

    embedding vector(768),```

    context_text TEXT,

    UNIQUE(chat_id, message_id)### Статистика PostgreSQL

);

``````sql

-- Общая статистика

## ТребованияSELECT 

    COUNT(*) as total_messages,

- Go 1.24+    COUNT(message_embedding) as vectorized_messages,

- PostgreSQL с расширением `pgvector`    ROUND(COUNT(message_embedding) * 100.0 / COUNT(*), 2) as vectorization_percent

- Gemini API ключFROM chat_messages;

- ~2GB RAM для больших архивов

-- Статистика по чатам

## ПроизводительностьSELECT 

    chat_id,

- ~240 запросов в минуту к Gemini API    COUNT(*) as messages,

- Батчевая обработка по 50 сообщений    COUNT(message_embedding) as vectorized,

- Автоматическое управление rate limiting    ROUND(COUNT(message_embedding) * 100.0 / COUNT(*), 2) as percent

- Кэширование эмбеддингов для ускоренияFROM chat_messages 

GROUP BY chat_id

## ПоддержкаORDER BY messages DESC;

```

При возникновении проблем:

1. Запустите `go run cmd/test_setup/main.go` для диагностики## Оптимизация производительности

2. Проверьте логи в файле `migration.log`

3. Убедитесь в корректности настроек в `.env` файле### Rate Limiting

- **Tier 1 Gemini**: 240 запросов/мин, 24,000/день
- **Адаптивное замедление** при достижении лимитов
- **Экспоненциальная задержка** при ошибках

### Кэширование

- **Локальный кэш** эмбеддингов (`./cache/embeddings/`)
- **Автосохранение** каждые 100 эмбеддингов
- **Дедупликация** одинаковых контекстов

### Пакетная обработка

- **Размер пакета**: 50-100 сообщений
- **Пауза между пакетами**: 1-60 секунд
- **Параллельная обработка**: отключена для экономии API

## Управление процессом

### Пауза/возобновление

```bash
# Graceful остановка
Ctrl+C

# Возобновление
./migrate_history.exe
```

### Мониторинг API квот

```bash
# В логах ищите строки:
grep "API:" migration.log
grep "Rate limit" migration.log
grep "Запросов в минуту" migration.log
```

## Решение проблем

### Ошибки подключения к PostgreSQL

```bash
# Проверка подключения
psql -h 46.19.69.26 -p 5432 -U gen_user -d default_db
```

### Превышение лимитов API

- Подождите сброса лимитов (обновляется каждую минуту/день)
- Уменьшите `EMBEDDING_REQUESTS_PER_MINUTE` в .env
- Переключитесь на резервный API ключ

### Проблемы с памятью

- Уменьшите `EMBEDDING_BATCH_SIZE`
- Увеличьте `EMBEDDING_BATCH_DELAY`  
- Перезапустите процесс

### Ошибки парсинга JSON

- Проверьте кодировку файлов (должна быть UTF-8)
- Убедитесь в валидности JSON структуры
- Исключите поврежденные файлы

## Результат миграции

После успешной миграции:

1. **В PostgreSQL**: все сообщения с эмбеддингами
2. **В кэше**: локальные копии эмбеддингов  
3. **В логах**: полная статистика процесса
4. **В migration_state.json**: состояние для возобновления

### Проверка результата

```sql
-- Топ-10 чатов по количеству сообщений
SELECT chat_id, COUNT(*) as messages
FROM chat_messages 
GROUP BY chat_id 
ORDER BY messages DESC 
LIMIT 10;

-- Проверка векторного поиска
SELECT message_text, 
       message_embedding <=> '[0.1,0.2,0.3,...]'::vector as distance
FROM chat_messages 
WHERE message_embedding IS NOT NULL
ORDER BY distance ASC 
LIMIT 5;
```

## Дальнейшие шаги

1. **Переключение бота** на PostgreSQL storage
2. **Тестирование векторного поиска** 
3. **Настройка индексов** для оптимизации
4. **Мониторинг производительности**

---

**Автор**: Henry-Case-dev  
**Дата**: Август 2025  
**Версия**: 1.0.0
