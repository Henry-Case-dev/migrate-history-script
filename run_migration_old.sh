#!/bin/bash

echo "🚀 Запуск миграции истории чата в PostgreSQL с векторизацией"
echo "============================================================"

# Проверяем существование директории с историей
if [ ! -d "../../История чата" ]; then
    echo "❌ Директория 'История чата' не найдена!"
    echo "   Убедитесь, что запускаете скрипт из корня проекта"
    exit 1
fi

# Выводим информацию о настройках
echo "📊 Конфигурация:"
echo "   PostgreSQL: $POSTGRESQL_HOST:$POSTGRESQL_PORT"
echo "   База данных: $POSTGRESQL_DBNAME"
echo "   Пользователь: $POSTGRESQL_USER"
echo "   Модель эмбеддингов: $GEMINI_EMBEDDING_MODEL_NAME"
echo ""

# Проверяем размер истории
history_size=$(du -sh "../../История чата" 2>/dev/null | cut -f1)
echo "📁 Размер истории чата: $history_size"
echo ""

# Подсчитываем количество JSON файлов
json_count=$(find "../../История чата" -name "*.json" | wc -l)
echo "📄 Найдено JSON файлов: $json_count"
echo ""

echo "⚠️  ВАЖНО:"
echo "   - Миграция займет время (примерно 1-3 дня для векторизации)"
echo "   - Стоимость векторизации: ~$10-15"
echo "   - Процесс можно остановить Ctrl+C и возобновить позже"
echo "   - Состояние сохраняется в migration_state.json"
echo ""

read -p "🤔 Продолжить миграцию? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Миграция отменена"
    exit 0
fi

echo ""
echo "🔄 Запуск миграции..."
echo "   Логи сохраняются в migration.log"
echo "   Прогресс отображается в реальном времени"
echo ""

# Запускаем миграцию с логированием
./migrate_history.exe 2>&1 | tee migration.log

# Проверяем результат
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Миграция завершена успешно!"
    echo "   Проверьте migration_state.json для статистики"
    echo "   Логи сохранены в migration.log"
else
    echo ""
    echo "❌ Миграция завершилась с ошибкой"
    echo "   Проверьте migration.log для деталей"
    echo "   Вы можете возобновить миграцию позже"
    exit 1
fi
