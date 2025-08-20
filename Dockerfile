# Сборка бинарника
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/server

# Запуск
FROM alpine:latest

# Установка ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates
WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

EXPOSE 8080

CMD ["./main"]