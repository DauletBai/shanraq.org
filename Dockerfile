# Этап 1: Сборка приложения
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server ./cmd/server

# Этап 2: Создание минимального образа для запуска
FROM alpine:latest

WORKDIR /app

# Копируем веб-файлы (шаблоны и статику)
COPY --from=builder /app/web ./web

# Копируем только исполняемый файл из этапа сборки
COPY --from=builder /app/server .

# Открываем порт, на котором будет работать наше приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["/app/server"]