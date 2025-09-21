# 🚀 Migrate History Script# Migrate History Script# 🚀 Миграция истории чата в PostgreSQL с векторизацией



Скрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.



## ⚡ Быстрый стартСкрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.## Описание



### 1. Клонирование и установка



```bash## Быстрый старт

git clone https://github.com/Henry-Case-dev/migrate-history-script.git

cd migrate-history-scriptgit clone https://github.com/Henry-Case-dev/migrate-history-script.git

go mod downloadcd migrate-history-script

```go mod download



### 2. Настройка окруженияЭтот инструмент выполняет миграцию истории чата из JSON файлов в PostgreSQL с автоматической векторизацией для поиска по семантическому сходству.



Создайте файл `.env` в корне проекта:



```env### 1. Клонирование проекта## Основные возможности

# PostgreSQL

POSTGRESQL_HOST=your_host

POSTGRESQL_PORT=5432

POSTGRESQL_USER=your_user```bash- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры

POSTGRESQL_PASSWORD=your_password

POSTGRESQL_DBNAME=your_databasegit clone <repository-url>- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  



# Gemini APIcd migrate-history-script- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных

GEMINI_API_KEY=your_gemini_api_key

GEMINI_EMBEDDING_MODEL_NAME=embedding-001```- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени



# Настройки векторизации (опционально)- ✅ **Rate limiting** - защита от превышения квот API

EMBEDDING_REQUESTS_PER_MINUTE=240

EMBEDDING_REQUESTS_PER_DAY=24000### 2. Установка зависимостей- ✅ **Кэширование эмбеддингов** - экономия API вызовов

EMBEDDING_REQUEST_DELAY=0.3s

EMBEDDING_BATCH_SIZE=100- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

EMBEDDING_CACHE_ENABLED=true

``````bash



### 3. Подготовка данныхgo mod download## Требования



Поместите JSON файлы экспорта Telegram в папку `data/`:```



```- PostgreSQL с установленным pgvector

data/

├── chat_export_1.json### 3. Настройка конфигурации- Действующий API ключ Gemini (tier 1)

├── chat_export_2.json

└── ...- Go 1.21+

```

Создайте файл `.env` в корне проекта:- ~500MB свободного места в БД

### 4. Запуск

- ~$10-15 для векторизации

```bash

# Проверка настроек```env

go run cmd/test_setup/main.go

# PostgreSQL## Настройка

# Запуск миграции

go run main.goPOSTGRESQL_HOST=your_host



# Или скомпилированной версиейPOSTGRESQL_PORT=5432### 1. Переменные окружения (.env)

go build -o migrate_history.exe

./migrate_history.exePOSTGRESQL_USER=your_user

```

POSTGRESQL_PASSWORD=your_password```env

## 📋 Основные возможности

POSTGRESQL_DBNAME=your_database# PostgreSQL (уже настроено)

- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры

- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  POSTGRESQL_HOST=46.19.69.26

- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных

- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени# Gemini APIPOSTGRESQL_PORT=5432

- ✅ **Rate limiting** - защита от превышения квот API

- ✅ **Кэширование эмбеддингов** - экономия API вызововGEMINI_API_KEY=your_gemini_api_keyPOSTGRESQL_USER=gen_user

- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

GEMINI_EMBEDDING_MODEL_NAME=embedding-001POSTGRESQL_PASSWORD=poshelnahuy0880

## 📋 Требования

POSTGRESQL_DBNAME=default_db

- **Go 1.21+** 

- **PostgreSQL** с установленным pgvector# Настройки производительности (опционально)

- **Gemini API ключ** (tier 1 рекомендуется)

- **~500MB** свободного места в БДEMBEDDING_REQUESTS_PER_MINUTE=240# Настройки векторизации

- **~$10-15** для векторизации большого архива

EMBEDDING_REQUESTS_PER_DAY=24000EMBEDDING_REQUESTS_PER_MINUTE=240

## 🔧 Продвинутые настройки

EMBEDDING_REQUEST_DELAY=300msEMBEDDING_REQUESTS_PER_DAY=24000

### Переменные окружения

```EMBEDDING_BATCH_SIZE=100

```env

# ПроизводительностьEMBEDDING_REQUEST_DELAY=0.3s

EMBEDDING_BATCH_SIZE=100

EMBEDDING_BATCH_DELAY=60s### 4. Подготовка данныхEMBEDDING_BATCH_DELAY=60s

EMBEDDING_REQUESTS_PER_MINUTE=240

EMBEDDING_REQUEST_DELAY=0.3sEMBEDDING_ADAPTIVE_THROTTLING=true



# КэшированиеПоместите JSON файлы экспорта Telegram в папку `data/`:EMBEDDING_SAFETY_MARGIN=0.8

EMBEDDING_CACHE_ENABLED=true

EMBEDDING_CACHE_DIR=./cache/embeddingsEMBEDDING_MAX_RETRIES=3



# Возобновление миграции```EMBEDDING_CACHE_ENABLED=true

MIGRATION_RESUME_ENABLED=true

MIGRATION_STATE_FILE=./migration_state.jsondata/EMBEDDING_CACHE_DIR=./cache/embeddings

```

├── chat_export_1.jsonMIGRATION_RESUME_ENABLED=true

### Контекстная векторизация

├── chat_export_2.jsonMIGRATION_STATE_FILE=./migration_state.json

Скрипт использует адаптивные контекстные окна:

└── ...

- **Базовое окно**: 5+5 сообщений для обычных сообщений

- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  ```# Gemini API (проверьте свой ключ)

- **Максимальное окно**: 25+25 для важных дискуссий

- **Фильтрация шума**: исключение коротких/повторяющихся сообщенийGEMINI_API_KEY=your_api_key_here



## 🚦 Управление процессом### 5. ЗапускGEMINI_EMBEDDING_MODEL_NAME=embedding-001



### Запуск и остановка```



```bash#### Проверка настроек:

# Graceful остановка

Ctrl+C```bash### 2. Структура директорий



# Возобновление (автоматически продолжит с места остановки)go run cmd/test_setup/main.go

./migrate_history.exe

``````

# Запуск в фоне

nohup ./migrate_history.exe > migration.log 2>&1 &rofloslav/

```

#### Запуск миграции:├── cmd/migrate_history/

### Мониторинг

```bash│   ├── migrate_history.exe     # Скомпилированный инструмент

```bash

# Логи в реальном времениgo run main.go│   ├── run_migration.sh       # Скрипт запуска

tail -f migration.log

```│   └── main.go               # Исходный код

# Статус миграции

cat migration_state.json | jq '.'└── История чата/             # JSON файлы истории



# API квотыили скомпилированной версией:    ├── chat_history_1.json

grep "Rate limit\|API:" migration.log

``````bash    ├── chat_history_2.json



## 🛡️ Безопасность при переносеgo build -o migrate_history.exe    └── ...



При переносе скрипта на другую машину **уже обработанные данные не будут загружаться повторно** благодаря:./migrate_history.exe```



1. **Отслеживание состояния** - файл `migration_state.json` хранит информацию о обработанных файлах```

2. **Database constraints** - PostgreSQL UNIQUE ограничения на `(chat_id, message_id)`

3. **ON CONFLICT стратегия** - `INSERT ... ON CONFLICT DO NOTHING` предотвращает дублирование## Использование



**Для переноса:**## Возможности

1. Скопируйте `migration_state.json` 

2. Убедитесь что PostgreSQL содержит уже импортированные данные### Быстрый старт

3. Запустите скрипт - он автоматически пропустит обработанные файлы

- ✅ **Импорт JSON → PostgreSQL** - перенос истории с сохранением структуры

## 📊 Структура базы данных

- ✅ **Контекстная векторизация** - умные эмбеддинги с учетом соседних сообщений  ```bash

Скрипт автоматически создает таблицу `messages`:

- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данныхcd cmd/migrate_history

```sql

CREATE TABLE messages (- ✅ **Защита от дублирования** - безопасный повторный запуск./run_migration.sh

    id SERIAL PRIMARY KEY,

    chat_id BIGINT NOT NULL,- ✅ **Rate limiting** - защита от превышения квот API```

    message_id INTEGER NOT NULL,

    user_id BIGINT,- ✅ **Кэширование эмбеддингов** - экономия API вызовов

    username TEXT,

    text TEXT,# В фоне:

    date_sent TIMESTAMP WITH TIME ZONE,

    embedding vector(768),## Структура проекта```bash

    context_text TEXT,

    UNIQUE(chat_id, message_id)nohup ./migrate_history.exe > migration.log 2>&1 &

);

`````````



### Проверка результатаmigrate-history-script/



```sql├── main.go                    # Основной файл программы### Ручной запуск

-- Общая статистика

SELECT ├── internal/                  # Внутренние пакеты

    COUNT(*) as total_messages,

    COUNT(embedding) as vectorized_messages,│   ├── config/               # Конфигурация```bash

    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as vectorization_percent

FROM messages;│   ├── storage/              # Работа с PostgreSQLcd cmd/migrate_history



-- Топ чаты по количеству сообщений│   ├── gemini/               # Клиент Gemini API./migrate_history.exe

SELECT chat_id, COUNT(*) as messages

FROM messages │   ├── llm/                  # Интерфейсы LLM```

GROUP BY chat_id 

ORDER BY messages DESC │   └── utils/                # Утилиты

LIMIT 10;

```├── cmd/### Возобновление прерванной миграции



## ⚡ Производительность│   └── test_setup/           # Скрипт проверки настроек



### Rate Limiting├── data/                     # JSON файлы экспорта TelegramПросто запустите инструмент снова - он автоматически продолжит с места остановки используя `migration_state.json`.

- **Tier 1 Gemini**: 240 запросов/мин, 24,000/день

- **Адаптивное замедление** при достижении лимитов├── cache/                    # Кэш эмбеддингов

- **Экспоненциальная задержка** при ошибках

├── migration_state.json     # Состояние миграции (автосоздается)## Стратегия векторизации

### Оптимизация

- **Пакетная обработка**: 50-100 сообщений за раз└── .env                      # Конфигурация (создать самостоятельно)

- **Локальный кэш** эмбеддингов для ускорения

- **Дедупликация** одинаковых контекстов```### Контекстные окна



## 🔧 Решение проблем



### Ошибки подключения к PostgreSQL## Управление процессом- **Базовое окно**: 5+5 сообщений для обычных сообщений

```bash

# Проверка подключения- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  

psql -h your_host -p 5432 -U your_user -d your_database

```- **Остановка**: `Ctrl+C` (состояние сохраняется автоматически)- **Максимальное окно**: 25+25 для важных дискуссий



### Превышение лимитов API- **Возобновление**: Перезапустите программу - автоматически продолжит с места остановки- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

- Подождите сброса лимитов (обновляется каждую минуту/день)

- Уменьшите `EMBEDDING_REQUESTS_PER_MINUTE` в .env- **Мониторинг**: Логи выводятся в консоль и файл `migration.log`

- Используйте резервный API ключ

### Адаптивная логика

### Проблемы с памятью

- Уменьшите `EMBEDDING_BATCH_SIZE`## Безопасность при переносе

- Увеличьте `EMBEDDING_BATCH_DELAY`  

- Перезапустите процесс```go



## 📁 Структура проектаПри переносе скрипта на другую машину:// Определение размера окна



```if isImportantMessage(msg) {        // Длинные сообщения, вопросы

migrate-history-script/

├── main.go                    # Основной файл программы1. Скопируйте файл `migration_state.json`     return 25                       // Максимальный контекст

├── internal/                  # Внутренние пакеты

│   ├── config/               # Конфигурация2. Убедитесь что PostgreSQL содержит уже импортированные данные} else if hasReplyChain(msg) {      // Ответы в цепочке

│   ├── storage/              # Работа с PostgreSQL

│   ├── gemini/               # Клиент Gemini API3. Скрипт автоматически пропустит уже обработанные файлы и сообщения    return 15                       // Средний контекст  

│   ├── llm/                  # Интерфейсы LLM

│   └── utils/                # Утилиты} else if containsQuestions(msg) {   // Вопросы

├── cmd/

│   └── test_setup/           # Скрипт проверки настроекМеханизм защиты:    return 10                       // Увеличенный контекст

├── data/                     # JSON файлы экспорта Telegram

├── cache/                    # Кэш эмбеддингов- Отслеживание обработанных файлов в `migration_state.json`} else {

├── migration_state.json     # Состояние миграции (автосоздается)

└── .env                      # Конфигурация (создать самостоятельно)- Уникальные constraint'ы в PostgreSQL: `UNIQUE(chat_id, message_id)`    return 5                        // Базовый контекст

```

- `INSERT ... ON CONFLICT DO NOTHING` стратегия}

## 🎯 После миграции

```

1. **Переключение бота** на PostgreSQL storage

2. **Тестирование векторного поиска** ## База данных

3. **Настройка индексов** для оптимизации

4. **Мониторинг производительности**## Мониторинг



---Скрипт автоматически создает таблицу `messages`:



**Автор**: Henry-Case-dev  ### Логи в реальном времени

**Репозиторий**: https://github.com/Henry-Case-dev/migrate-history-script  

**Версия**: 1.0.0```sql

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
