package main

import (
	"os"
	"strings"
)

func init() {
	// Автоматически определяем, запускается ли test_setup
	args := os.Args
	execPath, _ := os.Executable()

	// Проверяем аргументы командной строки или имя исполняемого файла
	isTestMode := false
	for _, arg := range args {
		if strings.Contains(arg, "test_setup") {
			isTestMode = true
			break
		}
	}

	// Если запускается через go run test_setup.go, активируем тестовый режим
	if strings.Contains(execPath, "test_setup") || len(args) > 1 && strings.Contains(args[1], "test_setup.go") {
		isTestMode = true
	}

	if isTestMode {
		// Переопределяем main для тестового режима
		testSetup()
		os.Exit(0)
	}
}
