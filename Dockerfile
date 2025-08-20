# Сборка бинарника
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/app

# Запуск
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
