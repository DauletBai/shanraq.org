# Этап 1: Общий этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем основной веб-сервер
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server ./cmd/server

# Этап 2: Финальный образ для основного приложения (минималистичный)
FROM alpine:latest AS app-release

WORKDIR /app

COPY --from=builder /app/web ./web
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["/app/server"]

# Этап 3: Финальный образ для seeder (с Go внутри)
FROM golang:1.24-alpine AS seeder-release

WORKDIR /app

COPY --from=builder /app .

# Команда по умолчанию не нужна, так как мы указываем ее в docker-compose