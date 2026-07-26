# --- Этап 1: Сборка (builder) ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Сначала копируем вообще всё (и код, и go.mod)
COPY . .

# Команда tidy сама найдет все зависимости в коде, скачает их и создаст go.sum
RUN go mod tidy

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# --- Этап 2: Финальный образ (runtime) ---
FROM alpine:3.19

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /server /app/server

USER appuser

EXPOSE 8080

CMD ["/app/server"]
