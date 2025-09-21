# 🚀 Миграция истории чата в PostgreSQL с векторизацией

## Описание

Этот инструмент выполняет миграцию истории чата из JSON файлов в PostgreSQL с автоматической векторизацией для поиска по семантическому сходству.

## Основные возможности

- ✅ **Миграция JSON → PostgreSQL** - перенос всей истории с сохранением структуры
- ✅ **Контекстная векторизация** - умное создание эмбеддингов с учетом контекста  
- ✅ **Возобновляемая миграция** - остановка/возобновление без потери данных
- ✅ **Прогресс-бар** - отслеживание процесса в реальном времени
- ✅ **Rate limiting** - защита от превышения квот API
- ✅ **Кэширование эмбеддингов** - экономия API вызовов
- ✅ **Фильтрация сообщений** - обработка только релевантных сообщений

## Требования

- PostgreSQL с установленным pgvector
- Действующий API ключ Gemini (tier 1)
- Go 1.21+
- ~500MB свободного места в БД
- ~$10-15 для векторизации

## Настройка

### 1. Переменные окружения (.env)

```env
# PostgreSQL (уже настроено)
POSTGRESQL_HOST=46.19.69.26
POSTGRESQL_PORT=5432
POSTGRESQL_USER=gen_user
POSTGRESQL_PASSWORD=poshelnahuy0880
POSTGRESQL_DBNAME=default_db

# Настройки векторизации
EMBEDDING_REQUESTS_PER_MINUTE=240
EMBEDDING_REQUESTS_PER_DAY=24000
EMBEDDING_BATCH_SIZE=100
EMBEDDING_REQUEST_DELAY=0.3s
EMBEDDING_BATCH_DELAY=60s
EMBEDDING_ADAPTIVE_THROTTLING=true
EMBEDDING_SAFETY_MARGIN=0.8
EMBEDDING_MAX_RETRIES=3
EMBEDDING_CACHE_ENABLED=true
EMBEDDING_CACHE_DIR=./cache/embeddings
MIGRATION_RESUME_ENABLED=true
MIGRATION_STATE_FILE=./migration_state.json

# Gemini API (проверьте свой ключ)
GEMINI_API_KEY=your_api_key_here
GEMINI_EMBEDDING_MODEL_NAME=embedding-001
```

### 2. Структура директорий

```
rofloslav/
├── cmd/migrate_history/
│   ├── migrate_history.exe     # Скомпилированный инструмент
│   ├── run_migration.sh       # Скрипт запуска
│   └── main.go               # Исходный код
└── История чата/             # JSON файлы истории
    ├── chat_history_1.json
    ├── chat_history_2.json
    └── ...
```

## Использование

### Быстрый старт

```bash
cd cmd/migrate_history
./run_migration.sh
```

# В фоне:
```bash
nohup ./migrate_history.exe > migration.log 2>&1 &
```

### Ручной запуск

```bash
cd cmd/migrate_history
./migrate_history.exe
```

### Возобновление прерванной миграции

Просто запустите инструмент снова - он автоматически продолжит с места остановки используя `migration_state.json`.

## Стратегия векторизации

### Контекстные окна

- **Базовое окно**: 5+5 сообщений для обычных сообщений
- **Расширенное окно**: 15+15 для связанных диалогов (reply chains)  
- **Максимальное окно**: 25+25 для важных дискуссий
- **Фильтрация шума**: исключение коротких/повторяющихся сообщений

### Адаптивная логика

```go
// Определение размера окна
if isImportantMessage(msg) {        // Длинные сообщения, вопросы
    return 25                       // Максимальный контекст
} else if hasReplyChain(msg) {      // Ответы в цепочке
    return 15                       // Средний контекст  
} else if containsQuestions(msg) {   // Вопросы
    return 10                       // Увеличенный контекст
} else {
    return 5                        // Базовый контекст
}
```

## Мониторинг

### Логи в реальном времени

```bash
tail -f migration.log
```

### Статус миграции

```bash
cat migration_state.json | jq '.'
```

### Статистика PostgreSQL

```sql
-- Общая статистика
SELECT 
    COUNT(*) as total_messages,
    COUNT(message_embedding) as vectorized_messages,
    ROUND(COUNT(message_embedding) * 100.0 / COUNT(*), 2) as vectorization_percent
FROM chat_messages;

-- Статистика по чатам
SELECT 
    chat_id,
    COUNT(*) as messages,
    COUNT(message_embedding) as vectorized,
    ROUND(COUNT(message_embedding) * 100.0 / COUNT(*), 2) as percent
FROM chat_messages 
GROUP BY chat_id
ORDER BY messages DESC;
```

## Оптимизация производительности

### Rate Limiting

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
