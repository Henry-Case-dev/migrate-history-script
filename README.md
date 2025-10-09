# 🚀 Telegram History Migration Tool# 🚀 Telegram History Migration Tool# 🚀 Migrate History Script



> Автоматическая миграция истории Telegram чатов в PostgreSQL с AI-векторизацией через Gemini API



[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)> Автоматическая миграция истории Telegram чатов в PostgreSQL с AI-векторизацией через Gemini APIСкрипт для миграции истории Telegram чатов в PostgreSQL с векторизацией сообщений через Gemini API.

[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)



---[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)## ⚡ Быстрый старт



## ⚡ Быстрый старт[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)



### 1. Установка Go (если еще не установлен)[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)### 1. Клонирование и установка



**Windows:**

```bash

# Скачайте установщик с официального сайта---```bash

https://go.dev/dl/

git clone https://github.com/Henry-Case-dev/migrate-history-script.git

# Или через Chocolatey

choco install golang## 📖 Содержаниеcd migrate-history-script

```

go mod download

**Linux/Mac:**

```bash- [Быстрый старт](#-быстрый-старт)```

# Linux

sudo apt-get install golang-go- [Требования](#-требования)



# Mac- [Основные возможности](#-основные-возможности)### 2. Настройка окружения

brew install go

```- [Конфигурация](#-конфигурация)



Проверьте установку:- [Режимы работы](#-режимы-работы)Создайте файл `.env` в корне проекта:

```bash

go version  # Должно показать Go 1.21+- [Управление процессом](#-управление-процессом)

```

- [Решение проблем](#-решение-проблем)```env

### 2. Клонирование и установка зависимостей

# PostgreSQL

```bash

# Клонируйте репозиторий---POSTGRESQL_HOST=your_host

git clone https://github.com/Henry-Case-dev/migrate-history-script.git

cd migrate-history-scriptPOSTGRESQL_PORT=5432



# Установите зависимости## ⚡ Быстрый стартPOSTGRESQL_USER=your_user

go mod download

```POSTGRESQL_PASSWORD=your_password



### 3. Настройка окружения### Шаг 1: УстановкаPOSTGRESQL_DBNAME=your_database



Создайте файл `.env` в корне проекта:



```env```bash# Gemini API

# PostgreSQL

POSTGRESQL_HOST=your_host# Клонируйте репозиторийGEMINI_API_KEY=your_gemini_api_key

POSTGRESQL_PORT=5432

POSTGRESQL_USER=your_usergit clone https://github.com/Henry-Case-dev/migrate-history-script.gitGEMINI_EMBEDDING_MODEL_NAME=embedding-001

POSTGRESQL_PASSWORD=your_password

POSTGRESQL_DBNAME=your_databasecd migrate-history-script



# Gemini API# Настройки векторизации (опционально)

GEMINI_API_KEY=your_gemini_api_key

GEMINI_EMBEDDING_MODEL_NAME=embedding-001# Установите зависимостиEMBEDDING_REQUESTS_PER_MINUTE=240



# Retry настройкиgo mod downloadEMBEDDING_REQUESTS_PER_DAY=24000

EMBEDDING_MAX_RETRIES=0  # 0 = бесконечный режим с 24ч таймаутом

```EMBEDDING_REQUEST_DELAY=0.3s

# Rate Limiting (рекомендуемые значения)

EMBEDDING_REQUESTS_PER_MINUTE=240EMBEDDING_BATCH_SIZE=100

EMBEDDING_REQUESTS_PER_DAY=24000

EMBEDDING_REQUEST_DELAY=0.3s### Шаг 2: Настройка `.env` файлаEMBEDDING_CACHE_ENABLED=true

```



### 4. Подготовка данных

Создайте файл `.env` в корне проекта:# Retry настройки

Поместите JSON файлы экспорта Telegram в папку `data/`:

EMBEDDING_MAX_RETRIES=5  # 0 = бесконечный режим с 24ч таймаутом

```

data/```env```

├── chat_export_1.json

├── chat_export_2.json# ============================================

└── ...

```# PostgreSQL Configuration### 3. Подготовка данных



### 5. Запуск# ============================================



```bashPOSTGRESQL_HOST=localhostПоместите JSON файлы экспорта Telegram в папку `data/`:

# Проверка настроек (опционально)

go run cmd/test_setup/main.goPOSTGRESQL_PORT=5432



# Запуск миграцииPOSTGRESQL_USER=your_user```

go run main.go

POSTGRESQL_PASSWORD=your_passworddata/

# Или скомпилируйте и запустите

go build -o migrate_history.exePOSTGRESQL_DBNAME=your_database├── chat_export_1.json

./migrate_history.exe

```├── chat_export_2.json



### 6. Остановка# ============================================└── ...



```bash# Gemini API Configuration```

# Graceful остановка (Ctrl+C)

# Программа сохранит состояние и можно будет продолжить позже# ============================================



# При следующем запуске программа автоматически продолжит с места остановкиGEMINI_API_KEY=your_gemini_api_key### 4. Запуск

./migrate_history.exe

```GEMINI_EMBEDDING_MODEL_NAME=embedding-001



---```bash



## 📋 Основные возможности# ============================================# Проверка настроек



- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры# Migration Settingsgo run cmd/test_setup/main.go

- ✅ **AI-векторизация** - умное создание эмбеддингов через Gemini API

- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных# ============================================

- ✅ **Бесконечный retry режим** - гарантия обработки всех сообщений (EMBEDDING_MAX_RETRIES=0)

- ✅ **Защита от дубликатов** - автоматический пропуск уже обработанных сообщенийEMBEDDING_MAX_RETRIES=5              # 0 = бесконечные retry (24ч таймаут)# Запуск миграции

- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени

- ✅ **Rate limiting** - защита от превышения квот APIEMBEDDING_BATCH_SIZE=100             # Размер пакета сообщенийgo run main.go

- ✅ **Кэширование эмбеддингов** - экономия API вызовов

EMBEDDING_CACHE_ENABLED=true        # Кэширование эмбеддингов

---

# Или скомпилированной версией

## 📋 Требования

# ============================================go build -o migrate_history.exe

| Компонент | Версия | Примечание |

|-----------|--------|------------|# Rate Limiting (рекомендуемые значения)./migrate_history.exe

| **Go** | 1.21+ | Язык программирования |

| **PostgreSQL** | 13+ с pgvector | База данных с векторным расширением |# ============================================```

| **Gemini API** | Tier 1+ | Для генерации эмбеддингов (240 req/min) |

| **Дисковое пространство** | ~500MB | Для хранения данных |EMBEDDING_REQUESTS_PER_MINUTE=240



### Установка pgvectorEMBEDDING_REQUESTS_PER_DAY=24000## 📋 Основные возможности



```sqlEMBEDDING_REQUEST_DELAY=0.3s

-- В PostgreSQL выполните:

CREATE EXTENSION IF NOT EXISTS vector;```- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры

```

- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  

---

### Шаг 3: Подготовка данных- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных

## 🔄 Режимы работы

- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени

### 🛡️ Бесконечный режим (РЕКОМЕНДУЕТСЯ для production)

Поместите JSON файлы экспорта Telegram в папку `data/`:- ✅ **Rate limiting** - защита от превышения квот API

```env

EMBEDDING_MAX_RETRIES=0  # Бесконечный режим- ✅ **Кэширование эмбеддингов** - экономия API вызовов

```

```- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

**Поведение:**

- ✅ **Никогда не пропускает сообщения** - гарантия полноты данныхdata/

- ✅ Повторяет попытки обработки до успеха или таймаута (24 часа)

- ✅ Экспоненциальная задержка до **60 минут** между попытками├── chat_export_1.json## 📋 Требования

- ✅ После таймаута - **завершение программы** с сохранением состояния

- ✅ При перезапуске продолжит с того же сообщения├── chat_export_2.json



**Когда использовать:** production, когда важна полнота данных└── chat_export_N.json- **Go 1.21+** 



### 🎯 Ограниченный режим```- **PostgreSQL** с установленным pgvector



```env- **Gemini API ключ** (tier 1 рекомендуется)

EMBEDDING_MAX_RETRIES=5  # Количество попыток

```### Шаг 4: Проверка и запуск- **~500MB** свободного места в БД



**Поведение:**- **~$10-15** для векторизации большого архива

- ✅ Делает 5 попыток при ошибках

- ✅ Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (max 60с)```bash

- ⚠️ При исчерпании попыток - возвращает ошибку

# Проверьте подключение к БД и API## 🔧 Продвинутые настройки

**Когда использовать:** тестирование, быстрая миграция

go run cmd/test_setup/main.go

---

### Переменные окружения

## 🚦 Управление процессом

# Запустите миграцию

### Запуск и остановка

go run main.go```env

```bash

# Запуск# Производительность

./migrate_history.exe

# Или скомпилируйте и запуститеEMBEDDING_BATCH_SIZE=100

# Graceful остановка (сохраняет состояние)

Ctrl+Cgo build -o migrate_history.exeEMBEDDING_BATCH_DELAY=60s



# Возобновление (автоматически продолжит)./migrate_history.exeEMBEDDING_REQUESTS_PER_MINUTE=240

./migrate_history.exe

```EMBEDDING_REQUEST_DELAY=0.3s

# Запуск в фоне (Linux/Mac)

nohup ./migrate_history.exe > migration.log 2>&1 &



# Windows (PowerShell)---# Retry механизм

Start-Process -NoNewWindow -FilePath ".\migrate_history.exe" -RedirectStandardOutput "migration.log" -RedirectStandardError "errors.log"

```EMBEDDING_MAX_RETRIES=5        # Количество попыток при ошибках



### Мониторинг## 📋 Требования                               # 0 = бесконечный режим с 24ч таймаутом



```bash

# Логи в реальном времени

tail -f migration.log| Компонент | Версия | Примечание |# Кэширование



# Статус миграции|-----------|--------|------------|EMBEDDING_CACHE_ENABLED=true

cat migration_state.json | jq '.'

| **Go** | 1.21+ | Язык программирования |EMBEDDING_CACHE_DIR=./cache/embeddings

# Проверка квот API

grep "Rate limit\|API:" migration.log| **PostgreSQL** | 13+ с pgvector | База данных с векторным расширением |

```

| **Gemini API** | Tier 1+ | Для генерации эмбеддингов |# Возобновление миграции

### Проверка результатов в БД

| **Дисковое пространство** | ~500MB | Для хранения данных |MIGRATION_RESUME_ENABLED=true

```sql

-- Общая статистика| **Бюджет API** | ~$10-15 | Для большого архива (ориентировочно) |MIGRATION_STATE_FILE=./migration_state.json

SELECT 

    COUNT(*) as total_messages,```

    COUNT(embedding) as vectorized_messages,

    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as percent### Установка pgvector

FROM messages;

### Retry механизм

-- Топ-10 чатов

SELECT chat_id, COUNT(*) as messages```sql

FROM messages 

GROUP BY chat_id -- В PostgreSQL выполните:Скрипт поддерживает два режима обработки ошибок:

ORDER BY messages DESC 

LIMIT 10;CREATE EXTENSION IF NOT EXISTS vector;

```

```#### Ограниченный режим (EMBEDDING_MAX_RETRIES > 0)

---

- Делает указанное количество попыток при ошибках

## 🛡️ Защита от дубликатов

---- Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (максимум 60с)

При переносе скрипта на другую машину **данные не дублируются**:

- При исчерпании попыток возвращает ошибку

1. ✅ `migration_state.json` - отслеживает обработанные файлы и индексы

2. ✅ PostgreSQL UNIQUE constraint на `(chat_id, message_id)`## ✨ Основные возможности

3. ✅ `INSERT ... ON CONFLICT DO NOTHING` стратегия

**Пример**: `EMBEDDING_MAX_RETRIES=5` - 5 попыток, затем ошибка

**Для переноса:**

```bash| Возможность | Описание |

# 1. Скопируйте файл состояния

scp migration_state.json user@server:/path/|------------|----------|#### Бесконечный режим (EMBEDDING_MAX_RETRIES = 0)



# 2. Убедитесь что PostgreSQL доступна| 🔄 **Миграция JSON → PostgreSQL** | Полный перенос истории с сохранением структуры |- **Никогда не пропускает сообщения** - гарантия полноты данных

psql -h your_host -U your_user -d your_database -c "SELECT COUNT(*) FROM messages;"

| 🧠 **AI-векторизация** | Контекстные эмбеддинги через Gemini API |- Повторяет попытки обработки до успеха или таймаута

# 3. Запустите - программа продолжит с места остановки

./migrate_history.exe| 💾 **Возобновляемость** | Остановка/возобновление без потери прогресса |- Экспоненциальная задержка до **60 минут** между попытками

```

| 📊 **Прогресс-бар** | Отслеживание процесса в реальном времени |- После достижения 60 минут - фиксированные интервалы в 60 минут

---

| ⚡ **Rate Limiting** | Автоматическая защита от превышения квот |- **Максимальная длительность**: 24 часа на одно сообщение

## 🔧 Решение проблем

| 💰 **Кэширование** | Экономия API вызовов |- При таймауте - **завершение программы** с сохранением состояния

### ❌ Ошибка: "Cannot connect to PostgreSQL"

| 🔍 **Фильтрация** | Обработка только релевантных сообщений |

```bash

# Проверьте подключение| 🛡️ **Защита от дубликатов** | UNIQUE constraints на message_id |**Рекомендуется для production**: гарантирует что все сообщения будут обработаны

psql -h your_host -p 5432 -U your_user -d your_database



# Проверьте расширение pgvector

psql -h your_host -U your_user -d your_database -c "SELECT * FROM pg_extension WHERE extname='vector';"---**Когда использовать:**

```

- `EMBEDDING_MAX_RETRIES=5` - для тестирования и быстрой миграции

**Решение:**

- Проверьте `.env` настройки (host, port, user, password)## ⚙️ Конфигурация- `EMBEDDING_MAX_RETRIES=0` - для production, когда важна полнота данных

- Убедитесь что PostgreSQL запущен

- Установите pgvector: `CREATE EXTENSION vector;`



### ⏱️ Таймаут 24 часа (EMBEDDING_MAX_RETRIES=0)### Основные параметры### Контекстная векторизация



```

❌ [CRITICAL] Не удалось обработать сообщение 12345 за 24 часа.

``````envСкрипт использует адаптивные контекстные окна:



**Причины:**# Производительность

- Проблемы с Gemini API (перегрузка сервиса)

- Исчерпание квот APIEMBEDDING_BATCH_SIZE=100              # Сообщений в пакете (50-200)- **Базовое окно**: 5+5 сообщений для обычных сообщений

- Нестабильное интернет-соединение

EMBEDDING_BATCH_DELAY=60s             # Задержка между пакетами- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  

**Решение:**

1. Проверьте статус Gemini API- **Максимальное окно**: 25+25 для важных дискуссий

2. Проверьте квоты: https://aistudio.google.com/app/apikey

3. Проверьте интернет-соединение# Rate Limiting- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

4. После устранения - просто перезапустите скрипт

EMBEDDING_REQUESTS_PER_MINUTE=240    # Лимит запросов/минуту

### 🚫 Превышение лимитов API

EMBEDDING_REQUESTS_PER_DAY=24000     # Лимит запросов/день## 🚦 Управление процессом

```bash

# Решение 1: Уменьшите лимиты в .envEMBEDDING_REQUEST_DELAY=0.3s         # Задержка между запросами

EMBEDDING_REQUESTS_PER_MINUTE=120  # Было 240

EMBEDDING_REQUEST_DELAY=0.6s       # Было 0.3s### Запуск и остановка



# Решение 2: Подождите сброса (автоматически)# Кэширование

# Tier 1: обновляется каждую минуту/день

EMBEDDING_CACHE_ENABLED=true         # Включить кэш```bash

# Решение 3: Используйте резервный ключ

GEMINI_API_KEY=your_backup_keyEMBEDDING_CACHE_DIR=./cache/embeddings# Graceful остановка

```

Ctrl+C

---

# Возобновление

## 📊 Структура базы данных

MIGRATION_RESUME_ENABLED=true        # Включить резюм# Возобновление (автоматически продолжит с места остановки)

```sql

CREATE TABLE messages (MIGRATION_STATE_FILE=./migration_state.json./migrate_history.exe

    id SERIAL PRIMARY KEY,

    chat_id BIGINT NOT NULL,```

    message_id INTEGER NOT NULL,

    user_id BIGINT,# Запуск в фоне

    username TEXT,

    text TEXT,### Контекстная векторизацияnohup ./migrate_history.exe > migration.log 2>&1 &

    date_sent TIMESTAMP WITH TIME ZONE,

    embedding vector(768),        -- Gemini embedding-001```

    context_text TEXT,

    UNIQUE(chat_id, message_id)Скрипт автоматически адаптирует размер контекста:

);

### Мониторинг

-- Индекс для векторного поиска

CREATE INDEX messages_embedding_idx - **Базовое окно**: 5±5 сообщений (обычные сообщения)

ON messages 

USING ivfflat (embedding vector_cosine_ops);- **Расширенное окно**: 15±15 (reply chains)```bash

```

- **Максимальное окно**: 25±25 (важные дискуссии)# Логи в реальном времени

---

- **Фильтрация шума**: исключение коротких/повторяющихся сообщенийtail -f migration.log

## 🎯 Векторный поиск (после миграции)



```sql

-- Поиск похожих сообщений (семантический поиск)---# Статус миграции

SELECT 

    text,cat migration_state.json | jq '.'

    embedding <=> '[your_vector]'::vector as distance

FROM messages ## 🔄 Режимы работы

WHERE embedding IS NOT NULL

ORDER BY distance ASC # API квоты

LIMIT 10;

### 🎯 Режим 1: Ограниченные retry (рекомендуется для тестирования)grep "Rate limit\|API:" migration.log

-- Поиск в конкретном чате

SELECT text, embedding <=> '[vector]'::vector as dist```

FROM messages 

WHERE chat_id = 123456 AND embedding IS NOT NULL```env

ORDER BY dist ASC 

LIMIT 5;EMBEDDING_MAX_RETRIES=5  # Количество попыток## 🛡️ Безопасность при переносе

```

```

---

При переносе скрипта на другую машину **уже обработанные данные не будут загружаться повторно** благодаря:

## 📁 Структура проекта

**Поведение:**

```

migrate-history-script/- ✅ Делает 5 попыток при ошибках1. **Отслеживание состояния** - файл `migration_state.json` хранит информацию о обработанных файлах

├── main.go                    # Основной файл (миграция, retry логика)

├── internal/- ✅ Экспоненциальная задержка: 2с → 4с → 8с → 16с → 32с (max 60с)2. **Database constraints** - PostgreSQL UNIQUE ограничения на `(chat_id, message_id)`

│   ├── config/                # Загрузка .env конфигурации

│   ├── storage/               # PostgreSQL клиент- ⚠️ При исчерпании попыток - возвращает ошибку3. **ON CONFLICT стратегия** - `INSERT ... ON CONFLICT DO NOTHING` предотвращает дублирование

│   ├── gemini/                # Gemini API клиент

│   ├── llm/                   # LLM интерфейсы

│   └── utils/                 # Утилиты

├── cmd/**Когда использовать:** быстрая миграция, тестирование**Для переноса:**

│   └── test_setup/            # Проверка настроек БД/API

├── data/                      # JSON экспорты Telegram (добавьте сюда)

├── cache/                     # Кэш эмбеддингов (автоматически)

├── migration_state.json       # Состояние (автоматически)---1. Скопируйте `migration_state.json` 

└── .env                       # Конфигурация (создайте вручную)

```2. Убедитесь что PostgreSQL содержит уже импортированные данные



---### 🛡️ Режим 2: Бесконечные retry (рекомендуется для production)3. Запустите скрипт - он автоматически пропустит обработанные файлы



## ⚡ Производительность



### Rate Limiting```env## 📊 Структура базы данных



| Tier | Requests/min | Requests/day | Рекомендации |EMBEDDING_MAX_RETRIES=0  # Бесконечный режим

|------|--------------|--------------|--------------|

| **Free** | 15 | 1,500 | Не подходит для больших архивов |```Скрипт автоматически создает таблицу `messages`:

| **Tier 1** | 240 | 24,000 | ✅ Рекомендуется |

| **Tier 2** | 1,000+ | 100,000+ | Для очень больших архивов |



### Оптимизация скорости**Поведение:**```sql



- ⚡ Пакетная обработка (50 сообщений в батче)- ✅ **Никогда не пропускает сообщения** - гарантия полноты данныхCREATE TABLE messages (

- 💾 Локальный кэш эмбеддингов

- 🔄 Дедупликация одинаковых контекстов- ✅ Повторяет до успеха или таймаута (24 часа)    id SERIAL PRIMARY KEY,

- 📉 Адаптивное замедление при лимитах

- 🎯 Возобновление с точного checkpoint без повторной обработки- ✅ Экспоненциальная задержка до 60 минут    chat_id BIGINT NOT NULL,



---- ✅ Затем фиксированные интервалы в 60 минут    message_id INTEGER NOT NULL,



## 📞 Поддержка- ✅ При таймауте - **завершение программы** с сохранением состояния    user_id BIGINT,



**Диагностика проблем:**    username TEXT,



```bash**Пример логов:**    text TEXT,

# 1. Запустите проверку настроек

go run cmd/test_setup/main.go    date_sent TIMESTAMP WITH TIME ZONE,



# 2. Проверьте логи```    embedding vector(768),

tail -n 100 migration.log

⏳ [Бесконечный режим] Попытка 15 | Прошло: 8h 30m | Осталось: 15h 30m | Следующая попытка через: 60m    context_text TEXT,

# 3. Проверьте состояние

cat migration_state.json | jq '.'✅ [Success] Сообщение 12345 успешно обработано после 15 попыток за 8h 35m    UNIQUE(chat_id, message_id)

```

```);

**Полезные ресурсы:**

- [Gemini API Documentation](https://ai.google.dev/docs)```

- [PostgreSQL pgvector](https://github.com/pgvector/pgvector)

- [Issues](https://github.com/Henry-Case-dev/migrate-history-script/issues)**Когда использовать:** production, когда важна полнота данных



---### Проверка результата



## 📄 Лицензия---



MIT License - см. [LICENSE](LICENSE)```sql



---## 🚦 Управление процессом-- Общая статистика



<div align="center">SELECT 



**Автор:** [Henry-Case-dev](https://github.com/Henry-Case-dev)  ### Запуск и остановка    COUNT(*) as total_messages,

**Репозиторий:** [migrate-history-script](https://github.com/Henry-Case-dev/migrate-history-script)  

    COUNT(embedding) as vectorized_messages,

⭐ Поставьте звезду, если проект был полезен!

```bash    ROUND(COUNT(embedding) * 100.0 / COUNT(*), 2) as vectorization_percent

</div>

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
