FROM golang:alpine AS builder

# Установить необходимые инструменты
RUN apk add --no-cache git curl

# Установить golang-migrate (конкретная версия совместимая с Go 1.23)
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:latest

# Установить необходимые утилиты
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

# Копировать бинарник migrate из builder
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Копировать приложение
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
