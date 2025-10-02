# Tenge-Web Dockerfile
# Агглютинативтік веб-қосымша Docker образы

# Build stage
FROM golang:1.21-alpine AS builder

# Жұмыс директориясын орнату
WORKDIR /app

# Go модульдерін көшіру
COPY go.mod go.sum ./
RUN go mod download

# Tenge компиляторын орнату
COPY --from=tenge-lang/tenge:latest /bin/tenge /usr/local/bin/tenge

# Кодды көшіру
COPY . .

# Tenge коддарын компиляциялау
RUN tenge compile ./backend/server/main.tng -o ./bin/shanraq-server
RUN tenge compile ./backend/api/main.tng -o ./bin/tenge-api-server
RUN tenge compile ./backend/orm/main.tng -o ./bin/tenge-orm

# Go коддарын құрастыру
RUN go build -o ./bin/shanraq-go ./backend/go/

# Production stage
FROM alpine:latest

# Жүйе пакеттерін орнату
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    nodejs \
    npm

# Жұмыс директориясын орнату
WORKDIR /app

# Бинарлық файлдарды көшіру
COPY --from=builder /app/bin/ /app/bin/

# Frontend құрастыру
COPY package*.json ./
RUN npm install
COPY frontend/ ./frontend/
RUN npm run build

# Статик файлдарды көшіру
COPY static/ ./static/

# Конфигурация файлдарын көшіру
COPY config/ ./config/

# Портты ашу
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# Серверді іске қосу
CMD ["./bin/shanraq-server"]

