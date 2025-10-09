# 🚀 Telegram History Migration Tool# 🚀 Migrate History Script



> Автоматическая миграция истории Telegram чатов в PostgreSQL с AI-векторизацией через Gemini APIСкрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.



[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)## ⚡ Быстрый старт

[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)### 1. Клонирование и установка



---```bash

git clone https://github.com/Henry-Case-dev/migrate-history-script.git

## 📖 Содержаниеcd migrate-history-script

go mod download

- [Быстрый старт](#-быстрый-старт)```

- [Требования](#-требования)

- [Основные возможности](#-основные-возможности)### 2. Настройка окружения

- [Конфигурация](#-конфигурация)

- [Режимы работы](#-режимы-работы)Создайте файл `.env` в корне проекта:

- [Управление процессом](#-управление-процессом)

- [Решение проблем](#-решение-проблем)```env

# PostgreSQL

---POSTGRESQL_HOST=your_host

POSTGRESQL_PORT=5432

## ⚡ Быстрый стартPOSTGRESQL_USER=your_user

POSTGRESQL_PASSWORD=your_password

### Шаг 1: УстановкаPOSTGRESQL_DBNAME=your_database



```bash# Gemini API

# Клонируйте репозиторийGEMINI_API_KEY=your_gemini_api_key

git clone https://github.com/Henry-Case-dev/migrate-history-script.gitGEMINI_EMBEDDING_MODEL_NAME=embedding-001

cd migrate-history-script

# Настройки векторизации (опционально)

# Установите зависимостиEMBEDDING_REQUESTS_PER_MINUTE=240

go mod downloadEMBEDDING_REQUESTS_PER_DAY=24000

```EMBEDDING_REQUEST_DELAY=0.3s

EMBEDDING_BATCH_SIZE=100

### Шаг 2: Настройка `.env` файлаEMBEDDING_CACHE_ENABLED=true



Создайте файл `.env` в корне проекта:# Retry настройки

EMBEDDING_MAX_RETRIES=5  # 0 = бесконечный режим с 24ч таймаутом

```env```

# ============================================

# PostgreSQL Configuration### 3. Подготовка данных

# ============================================

POSTGRESQL_HOST=localhostПоместите JSON файлы экспорта Telegram в папку `data/`:

POSTGRESQL_PORT=5432

POSTGRESQL_USER=your_user```

POSTGRESQL_PASSWORD=your_passworddata/

POSTGRESQL_DBNAME=your_database├── chat_export_1.json

├── chat_export_2.json

# ============================================└── ...

# Gemini API Configuration```

# ============================================

GEMINI_API_KEY=your_gemini_api_key### 4. Запуск

GEMINI_EMBEDDING_MODEL_NAME=embedding-001

```bash

# ============================================# Проверка настроек

# Migration Settingsgo run cmd/test_setup/main.go

# ============================================

EMBEDDING_MAX_RETRIES=5              # 0 = бесконечные retry (24ч таймаут)# Запуск миграции

EMBEDDING_BATCH_SIZE=100             # Размер пакета сообщенийgo run main.go

EMBEDDING_CACHE_ENABLED=true        # Кэширование эмбеддингов

# Или скомпилированной версией

# ============================================go build -o migrate_history.exe

# Rate Limiting (рекомендуемые значения)./migrate_history.exe

# ============================================```

EMBEDDING_REQUESTS_PER_MINUTE=240

EMBEDDING_REQUESTS_PER_DAY=24000## 📋 Основные возможности

EMBEDDING_REQUEST_DELAY=0.3s

```- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры

- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  

### Шаг 3: Подготовка данных- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных

- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени

Поместите JSON файлы экспорта Telegram в папку `data/`:- ✅ **Rate limiting** - защита от превышения квот API

- ✅ **Кэширование эмбеддингов** - экономия API вызовов

```- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

data/

├── chat_export_1.json## 📋 Требования

├── chat_export_2.json

└── chat_export_N.json- **Go 1.21+** 

```- **PostgreSQL** с установленным pgvector

- **Gemini API ключ** (tier 1 рекомендуется)

### Шаг 4: Проверка и запуск- **~500MB** свободного места в БД

- **~$10-15** для векторизации большого архива

```bash

# Проверьте подключение к БД и API## 🔧 Продвинутые настройки

go run cmd/test_setup/main.go

### Переменные окружения

# Запустите миграцию

go run main.go```env

# Производительность

# Или скомпилируйте и запуститеEMBEDDING_BATCH_SIZE=100

go build -o migrate_history.exeEMBEDDING_BATCH_DELAY=60s

./migrate_history.exeEMBEDDING_REQUESTS_PER_MINUTE=240

```EMBEDDING_REQUEST_DELAY=0.3s



---# Retry механизм

EMBEDDING_MAX_RETRIES=5        # Количество попыток при ошибках

## 📋 Требования                               # 0 = бесконечный режим с 24ч таймаутом



| Компонент | Версия | Примечание |# Кэширование

|-----------|--------|------------|EMBEDDING_CACHE_ENABLED=true

| **Go** | 1.21+ | Язык программирования |EMBEDDING_CACHE_DIR=./cache/embeddings

| **PostgreSQL** | 13+ с pgvector | База данных с векторным расширением |

| **Gemini API** | Tier 1+ | Для генерации эмбеддингов |# Возобновление миграции

| **Дисковое пространство** | ~500MB | Для хранения данных |MIGRATION_RESUME_ENABLED=true

| **Бюджет API** | ~$10-15 | Для большого архива (ориентировочно) |MIGRATION_STATE_FILE=./migration_state.json

```

### Установка pgvector

### Retry механизм

```sql

-- В PostgreSQL выполните:Скрипт поддерживает два режима обработки ошибок:

CREATE EXTENSION IF NOT EXISTS vector;

```#### Ограниченный режим (EMBEDDING_MAX_RETRIES > 0)

- Делает указанное количество попыток при ошибках

---- Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (максимум 60с)

- При исчерпании попыток возвращает ошибку

## ✨ Основные возможности

**Пример**: `EMBEDDING_MAX_RETRIES=5` - 5 попыток, затем ошибка

| Возможность | Описание |

|------------|----------|#### Бесконечный режим (EMBEDDING_MAX_RETRIES = 0)

| 🔄 **Миграция JSON → PostgreSQL** | Полный перенос истории с сохранением структуры |- **Никогда не пропускает сообщения** - гарантия полноты данных

| 🧠 **AI-векторизация** | Контекстные эмбеддинги через Gemini API |- Повторяет попытки обработки до успеха или таймаута

| 💾 **Возобновляемость** | Остановка/возобновление без потери прогресса |- Экспоненциальная задержка до **60 минут** между попытками

| 📊 **Прогресс-бар** | Отслеживание процесса в реальном времени |- После достижения 60 минут - фиксированные интервалы в 60 минут

| ⚡ **Rate Limiting** | Автоматическая защита от превышения квот |- **Максимальная длительность**: 24 часа на одно сообщение

| 💰 **Кэширование** | Экономия API вызовов |- При таймауте - **завершение программы** с сохранением состояния

| 🔍 **Фильтрация** | Обработка только релевантных сообщений |

| 🛡️ **Защита от дубликатов** | UNIQUE constraints на message_id |**Рекомендуется для production**: гарантирует что все сообщения будут обработаны



---**Когда использовать:**

- `EMBEDDING_MAX_RETRIES=5` - для тестирования и быстрой миграции

## ⚙️ Конфигурация- `EMBEDDING_MAX_RETRIES=0` - для production, когда важна полнота данных



### Основные параметры### Контекстная векторизация



```envСкрипт использует адаптивные контекстные окна:

# Производительность

EMBEDDING_BATCH_SIZE=100              # Сообщений в пакете (50-200)- **Базовое окно**: 5+5 сообщений для обычных сообщений

EMBEDDING_BATCH_DELAY=60s             # Задержка между пакетами- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  

- **Максимальное окно**: 25+25 для важных дискуссий

# Rate Limiting- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

EMBEDDING_REQUESTS_PER_MINUTE=240    # Лимит запросов/минуту

EMBEDDING_REQUESTS_PER_DAY=24000     # Лимит запросов/день## 🚦 Управление процессом

EMBEDDING_REQUEST_DELAY=0.3s         # Задержка между запросами

### Запуск и остановка

# Кэширование

EMBEDDING_CACHE_ENABLED=true         # Включить кэш```bash

EMBEDDING_CACHE_DIR=./cache/embeddings# Graceful остановка

Ctrl+C

# Возобновление

MIGRATION_RESUME_ENABLED=true        # Включить резюм# Возобновление (автоматически продолжит с места остановки)

MIGRATION_STATE_FILE=./migration_state.json./migrate_history.exe

```

# Запуск в фоне

### Контекстная векторизацияnohup ./migrate_history.exe > migration.log 2>&1 &

```

Скрипт автоматически адаптирует размер контекста:

### Мониторинг

- **Базовое окно**: 5±5 сообщений (обычные сообщения)

- **Расширенное окно**: 15±15 (reply chains)```bash

- **Максимальное окно**: 25±25 (важные дискуссии)# Логи в реальном времени

- **Фильтрация шума**: исключение коротких/повторяющихся сообщенийtail -f migration.log



---# Статус миграции

cat migration_state.json | jq '.'

## 🔄 Режимы работы

# API квоты

### 🎯 Режим 1: Ограниченные retry (рекомендуется для тестирования)grep "Rate limit\|API:" migration.log

```

```env

EMBEDDING_MAX_RETRIES=5  # Количество попыток## 🛡️ Безопасность при переносе

```

При переносе скрипта на другую машину **уже обработанные данные не будут загружаться повторно** благодаря:

**Поведение:**

- ✅ Делает 5 попыток при ошибках1. **Отслеживание состояния** - файл `migration_state.json` хранит информацию о обработанных файлах

- ✅ Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (max 60с)2. **Database constraints** - PostgreSQL UNIQUE ограничения на `(chat_id, message_id)`

- ⚠️ При исчерпании попыток - возвращает ошибку3. **ON CONFLICT стратегия** - `INSERT ... ON CONFLICT DO NOTHING` предотвращает дублирование



**Когда использовать:** быстрая миграция, тестирование**Для переноса:**



---1. Скопируйте `migration_state.json` 

2. Убедитесь что PostgreSQL содержит уже импортированные данные

### 🛡️ Режим 2: Бесконечные retry (рекомендуется для production)3. Запустите скрипт - он автоматически пропустит обработанные файлы



```env## 📊 Структура базы данных

EMBEDDING_MAX_RETRIES=0  # Бесконечный режим

```Скрипт автоматически создает таблицу `messages`:



**Поведение:**```sql

- ✅ **Никогда не пропускает сообщения** - гарантия полноты данныхCREATE TABLE messages (

- ✅ Повторяет до успеха или таймаута (24 часа)    id SERIAL PRIMARY KEY,

- ✅ Экспоненциальная задержка до 60 минут    chat_id BIGINT NOT NULL,

- ✅ Затем фиксированные интервалы в 60 минут    message_id INTEGER NOT NULL,

- ✅ При таймауте - **завершение программы** с сохранением состояния    user_id BIGINT,

    username TEXT,

**Пример логов:**    text TEXT,

    date_sent TIMESTAMP WITH TIME ZONE,

```    embedding vector(768),

⏳ [Бесконечный режим] Попытка 15 | Прошло: 8h 30m | Осталось: 15h 30m | Следующая попытка через: 60m    context_text TEXT,

✅ [Success] Сообщение 12345 успешно обработано после 15 попыток за 8h 35m    UNIQUE(chat_id, message_id)

```);

```

**Когда использовать:** production, когда важна полнота данных

### Проверка результата

---

```sql

## 🚦 Управление процессом-- Общая статистика

SELECT 

### Запуск и остановка    COUNT(*) as total_messages,

    COUNT(embedding) as vectorized_messages,

```bash    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as vectorization_percent

# ЗапускFROM messages;

./migrate_history.exe

-- Топ чаты по количеству сообщений

# Graceful остановка (сохраняет состояние)SELECT chat_id, COUNT(*) as messages

Ctrl+CFROM messages 

GROUP BY chat_id 

# Возобновление (автоматически продолжит)ORDER BY messages DESC 

./migrate_history.exeLIMIT 10;

```

# Запуск в фоне (Linux/Mac)

nohup ./migrate_history.exe > migration.log 2>&1 &## ⚡ Производительность

```

### Rate Limiting

### Мониторинг

- **Tier 1 Gemini**: 240 запросов/мин, 24,000/день

```bash- **Адаптивное замедление** при достижении лимитов

# Логи в реальном времени- **Экспоненциальная задержка** при ошибках

tail -f migration.log

### Оптимизация

# Статус миграции (JSON)

cat migration_state.json- **Пакетная обработка**: 50-100 сообщений за раз

- **Локальный кэш** эмбеддингов для ускорения

# Статистика (с jq)- **Дедупликация** одинаковых контекстов

cat migration_state.json | jq '.processed_count, .vectorized_count, .total_messages'

## 🔧 Решение проблем

# Проверка квот API

grep "Rate limit\|API:" migration.log### Ошибки подключения к PostgreSQL

```

```bash

### Проверка результатов в БД# Проверка подключения

psql -h your_host -p 5432 -U your_user -d your_database

```sql```

-- Общая статистика

SELECT ### Превышение лимитов API

    COUNT(*) as total_messages,

    COUNT(embedding) as vectorized_messages,- Подождите сброса лимитов (обновляется каждую минуту/день)

    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as percent- Уменьшите `EMBEDDING_REQUESTS_PER_MINUTE` в .env

FROM messages;- Используйте резервный API ключ



-- Топ-10 чатов### 24-часовой таймаут (EMBEDDING_MAX_RETRIES=0)

SELECT chat_id, COUNT(*) as messages

FROM messages Если программа завершилась с ошибкой "Не удалось обработать сообщение за 24 часа":

GROUP BY chat_id 

ORDER BY messages DESC 1. **Проверьте Gemini API** - возможно проблемы на стороне сервиса

LIMIT 10;2. **Проверьте API ключ** - убедитесь что квоты не исчерпаны

```3. **Проверьте подключение** - стабильность интернет соединения

4. **После устранения проблемы** просто перезапустите скрипт - он продолжит с места остановки

---

### Проблемы с памятью

## 🛡️ Безопасность и перенос

- Уменьшите `EMBEDDING_BATCH_SIZE`

### Защита от дубликатов- Увеличьте `EMBEDDING_BATCH_DELAY`  

- Перезапустите процесс

При переносе скрипта на другую машину **данные не дублируются**:

### Ошибки парсинга JSON

1. ✅ `migration_state.json` - отслеживание обработанных файлов

2. ✅ PostgreSQL UNIQUE constraint на `(chat_id, message_id)`- Проверьте кодировку файлов (должна быть UTF-8)

3. ✅ `INSERT ... ON CONFLICT DO NOTHING` стратегия- Убедитесь в валидности JSON структуры

- Исключите поврежденные файлы

### Как перенести на другой сервер

## 📁 Структура проекта

```bash

# 1. Скопируйте файлы```

scp migration_state.json user@server:/path/migrate-history-script/

scp .env user@server:/path/├── main.go                    # Основной файл программы

├── internal/                  # Внутренние пакеты

# 2. Убедитесь что PostgreSQL доступна│   ├── config/               # Конфигурация

psql -h your_host -U your_user -d your_database -c "SELECT COUNT(*) FROM messages;"│   ├── storage/              # Работа с PostgreSQL

│   ├── gemini/               # Клиент Gemini API

# 3. Запустите - он продолжит с места остановки│   ├── llm/                  # Интерфейсы LLM

./migrate_history.exe│   └── utils/                # Утилиты

```├── cmd/

│   └── test_setup/           # Скрипт проверки настроек

---├── data/                     # JSON файлы экспорта Telegram

├── cache/                    # Кэш эмбеддингов

## 🔧 Решение проблем├── migration_state.json     # Состояние миграции (автосоздается)

└── .env                      # Конфигурация (создать самостоятельно)

### ❌ Ошибка: "Cannot connect to PostgreSQL"```



```bash## 🎯 После миграции

# Проверьте подключение

psql -h your_host -p 5432 -U your_user -d your_database1. **Переключение бота** на PostgreSQL storage

2. **Тестирование векторного поиска** 

# Проверьте расширение pgvector3. **Настройка индексов** для оптимизации

psql -h your_host -U your_user -d your_database -c "SELECT * FROM pg_extension WHERE extname='vector';"4. **Мониторинг производительности**

```

## 📈 Векторный поиск

**Решение:**

- Проверьте `.env` настройки (host, port, user, password)После миграции вы сможете использовать семантический поиск:

- Убедитесь что PostgreSQL запущен

- Установите pgvector: `CREATE EXTENSION vector;````sql

-- Поиск похожих сообщений

---SELECT 

    text,

### ⏱️ Ошибка: "Превышено 24 часа" (EMBEDDING_MAX_RETRIES=0)    embedding <=> '[vector_here]'::vector as distance

FROM messages 

```WHERE embedding IS NOT NULL

❌ [CRITICAL] Не удалось обработать сообщение 12345 за 24 часа.ORDER BY distance ASC 

```LIMIT 10;



**Причины:**-- Создание индекса для быстрого поиска

- Проблемы с Gemini API (перегрузка сервиса)CREATE INDEX ON messages USING ivfflat (embedding vector_cosine_ops);

- Исчерпание квот API```

- Нестабильное интернет-соединение

## 📞 Поддержка

**Решение:**

1. Проверьте статус Gemini APIПри возникновении проблем:

2. Проверьте квоты: https://aistudio.google.com/app/apikey

3. Проверьте интернет-соединение1. Запустите `go run cmd/test_setup/main.go` для диагностики

4. После устранения - просто перезапустите скрипт2. Проверьте логи в файле `migration.log`

3. Убедитесь в корректности настроек в `.env` файле

---

---

### 🚫 Ошибка: "Rate limit exceeded"

**Автор**: Henry-Case-dev  

```bash**Репозиторий**: https://github.com/Henry-Case-dev/migrate-history-script  

# Решение 1: Уменьшите лимиты в .env**Версия**: 1.0.0  

EMBEDDING_REQUESTS_PER_MINUTE=120  # Было 240**Последнее обновление**: Сентябрь 2025
EMBEDDING_REQUEST_DELAY=0.6s       # Было 0.3s

# Решение 2: Подождите сброса (автоматически)
# Tier 1: обновляется каждую минуту/день

# Решение 3: Используйте резервный ключ
GEMINI_API_KEY=your_backup_key
```

---

### 💾 Проблемы с памятью

**Симптомы:** процесс убивается ОС, падает производительность

**Решение:**
```env
EMBEDDING_BATCH_SIZE=50              # Уменьшить с 100
EMBEDDING_BATCH_DELAY=120s           # Увеличить с 60s
```

---

### 📝 Ошибка парсинга JSON

**Решение:**
- Убедитесь в кодировке UTF-8
- Проверьте валидность JSON: `jq '.' data/your_file.json`
- Исключите поврежденные файлы из папки `data/`

---

## 📊 Структура базы данных

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

## 🎯 Векторный поиск (после миграции)

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

## 📁 Структура проекта

```
migrate-history-script/
├── 📄 main.go                    # Основной файл (миграция, retry логика)
├── 📦 internal/
│   ├── config/                  # Загрузка .env конфигурации
│   ├── storage/                 # PostgreSQL клиент
│   ├── gemini/                  # Gemini API клиент
│   ├── llm/                     # LLM интерфейсы
│   └── utils/                   # Утилиты
├── 🧪 cmd/
│   └── test_setup/              # Проверка настроек БД/API
├── 📂 data/                     # JSON экспорты Telegram (добавьте сюда)
├── 💾 cache/                    # Кэш эмбеддингов (автоматически)
├── 📊 migration_state.json      # Состояние (автоматически)
└── ⚙️ .env                       # Конфигурация (создайте вручную)
```

---

## ⚡ Производительность

### Rate Limiting

| Tier | Requests/min | Requests/day | Рекомендации |
|------|--------------|--------------|--------------|
| **Free** | 15 | 1,500 | Не подходит для больших архивов |
| **Tier 1** | 240 | 24,000 | ✅ Рекомендуется |
| **Tier 2** | 1,000+ | 100,000+ | Для очень больших архивов |

### Оптимизация скорости

- ⚡ Пакетная обработка (50-100 сообщений)
- 💾 Локальный кэш эмбеддингов
- 🔄 Дедупликация одинаковых контекстов
- 📉 Адаптивное замедление при лимитах

---

## 📞 Поддержка

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

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE)

---

<div align="center">

**Автор:** [Henry-Case-dev](https://github.com/Henry-Case-dev)  
**Репозиторий:** [migrate-history-script](https://github.com/Henry-Case-dev/migrate-history-script)  

⭐ Поставьте звезду, если проект был полезен!

</div>
