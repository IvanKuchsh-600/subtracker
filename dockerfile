FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Установка необходимых инструментов
RUN apk add --no-cache git ca-certificates

# Копирование go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subtracker ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копирование бинарника и конфигов
COPY --from=builder /app/subtracker .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./subtracker"]