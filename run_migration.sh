#!/bin/bash

echo "��� Luna Bot - Миграция истории чата"
echo "===================================="

go version
echo "✅ Go найден"

if [ ! -f ".env" ]; then
    echo "❌ .env не найден"
    exit 1
fi
echo "✅ .env найден"

mkdir -p data cache/embeddings tools
echo "✅ Директории созданы"

echo ""
echo "1) Тест настроек"
echo "2) Сборка"
echo "3) Запуск"
echo "4) Сборка + запуск"
echo ""

read -p "Выбор [1-4]: " choice

case $choice in
    1) go run test_setup.go ;;
    2) go build -o migrate_history.exe main.go && echo "✅ Собрано" ;;
    3) go run main.go ;;
    4) go build -o migrate_history.exe main.go && ./migrate_history.exe ;;
esac
