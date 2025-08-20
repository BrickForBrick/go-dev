.PHONY: help build run test clean docker-build docker-up docker-down swagger deps

# Переменные
APP_NAME=subscription-service
DOCKER_COMPOSE=docker-compose
GO_FILES=$(shell find . -name "*.go" -type f)

# Помощь
help: ## Показать эту справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Установка зависимостей
deps: ## Установить зависимости
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

# Генерация Swagger документации
swagger: ## Сгенерировать Swagger документацию
	swag init -g cmd/server/main.go -o ./docs

# Сборка приложения
build: swagger ## Собрать приложение
	go build -o bin/$(APP_NAME) ./cmd/server

# Запуск приложения локально
run: swagger ## Запустить приложение локально
	go run ./cmd/server

# Тесты
test: ## Запустить тесты
	go test -v ./...

# Линтинг
lint: ## Запустить линтер
	golangci-lint run

# Форматирование кода
fmt: ## Форматировать код
	go fmt ./...

# Очистка
clean: ## Удалить собранные файлы
	rm -rf bin/
	rm -rf docs/

# Docker команды
docker-build: ## Собрать Docker образ
	docker build -t $(APP_NAME) .

docker-up: ## Запустить сервисы через Docker Compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Остановить сервисы Docker Compose
	$(DOCKER_COMPOSE) down

docker-logs: ## Показать логи Docker контейнеров
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Перезапустить сервисы
	$(DOCKER_COMPOSE) restart

# База данных
db-migrate-up: ## Применить миграции
	go run ./cmd/server -migrate-up

db-migrate-down: ## Откатить миграции  
	go run ./cmd/server -migrate-down

# Полная пересборка и запуск
rebuild: clean build ## Полная пересборка

# Разработка
dev: swagger ## Запуск в режиме разработки
	air || go run ./cmd/server

# Проверка готовности к production
check: swagger lint test ## Проверки перед деплоем