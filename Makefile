.PHONY: build run docker-up docker-down clean generate migrate-up migrate-down install test

# Сборка бинарника
build:
	go build -o bin/server ./cmd/server

# Запуск приложения
run:
	go run ./cmd/server

# Запуск PostgreSQL
docker-up:
	docker-compose up -d postgres
	@echo "PostgreSQL running on port 5433"

# Остановка PostgreSQL
docker-down:
	docker-compose down

# Очистка
clean:
	rm -rf bin/
	docker-compose down -v

# Генерация кода sqlc
generate:
	sqlc generate

# Применение миграций
migrate-up:
	goose -dir migrations postgres "host=localhost port=5433 user=postgres password=postgres dbname=phoneservice sslmode=disable" up

# Откат миграций
migrate-down:
	goose -dir migrations postgres "host=localhost port=5433 user=postgres password=postgres dbname=phoneservice sslmode=disable" down

# Установка инструментов
install:
	go mod download
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Тестирование API
test:
	@echo "=== Health check ==="
	@curl -s http://localhost:8080/health | jq '.'
	@echo ""
	@echo "=== Import test ==="
	@curl -s -X POST http://localhost:8080/api/numbers/import \
		-H "Content-Type: application/json" \
		-d '{"numbers":["+79161234567","9123456789"],"source":"telegram"}' | jq '.'
	@echo ""
	@echo "=== Search test ==="
	@curl -s "http://localhost:8080/api/numbers/search?number=916" | jq '.'
