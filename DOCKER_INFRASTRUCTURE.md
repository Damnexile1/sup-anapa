# SUP-Anapa - Обновленная инфраструктура

## ✅ Что изменено

### Docker-first подход
- ✅ Все инструменты (golang-migrate) устанавливаются внутри Docker контейнера
- ✅ База данных доступна только внутри Docker сети (sup-anapa-network)
- ✅ Используется pgxpool для подключения к PostgreSQL
- ✅ Миграции выполняются внутри контейнера приложения
- ✅ Полная изоляция - ничего не устанавливается на хост-машину

### Обновленные репозитории
- ✅ Все репозитории переписаны на pgxpool
- ✅ Добавлен context.Context во все методы
- ✅ BookingRepository, InstructorRepository, SlotRepository, AdminRepository

### Makefile
- ✅ Убрана команда `install-tools` (всё в Docker)
- ✅ Все команды работают через `docker-compose exec`
- ✅ Миграции выполняются внутри контейнера
- ✅ Тесты и сборка внутри контейнера

## 🚀 Быстрый старт

```bash
# 1. Создать .env файл (если нужно)
make setup

# 2. Запустить всё одной командой
make deploy-local

# Приложение доступно на http://localhost:8080
```

## 📋 Основные команды

```bash
make help              # Показать все команды
make deploy-local      # Локальный деплой (Docker + миграции)
make deploy-prod       # Продакшн деплой
make docker-logs       # Логи всех контейнеров
make docker-logs-app   # Логи приложения
make docker-logs-db    # Логи БД
make migrate-up        # Применить миграции (внутри контейнера)
make db-shell          # Подключиться к БД
make shell             # Shell в контейнере приложения
```

## 🔧 Архитектура

```
┌─────────────────────────────────────┐
│         Host Machine                │
│                                     │
│  ┌───────────────────────────────┐ │
│  │   Docker Network              │ │
│  │   (sup-anapa-network)         │ │
│  │                               │ │
│  │  ┌──────────┐  ┌───────────┐ │ │
│  │  │   App    │  │ PostgreSQL│ │ │
│  │  │Container │◄─┤ Container │ │ │
│  │  │          │  │           │ │ │
│  │  │ :8080    │  │ :5432     │ │ │
│  │  └────┬─────┘  └───────────┘ │ │
│  │       │                       │ │
│  └───────┼───────────────────────┘ │
│          │                         │
│    Port Forward                    │
│      8080:8080                     │
└──────────┼─────────────────────────┘
           │
      localhost:8080
```

### Особенности:
- PostgreSQL доступен только внутри Docker сети
- Приложение подключается к БД через hostname `postgres`
- Порт 8080 пробрасывается на хост для доступа к приложению
- Миграции выполняются из контейнера приложения

## 📦 Зависимости

### Go модули
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL connection pool
- `golang.org/x/crypto/bcrypt` - Хеширование паролей (TODO)
- `github.com/gorilla/sessions` - Сессии (TODO)

### Docker образы
- `golang:1.21-alpine` - Builder
- `alpine:latest` - Runtime
- `postgres:15-alpine` - База данных

## 🗄️ База данных

### Подключение
```go
// Внутри контейнера
DATABASE_URL=postgres://postgres:postgres@postgres:5432/sup_anapa?sslmode=disable
```

### Миграции
```bash
# Применить все миграции
make migrate-up

# Откатить последнюю
make migrate-down

# Пересоздать всё
make migrate-reset

# Создать новую миграцию
make migrate-create NAME=add_reviews
```

## 🔄 Workflow разработки

### Первый запуск
```bash
make setup          # Создать .env
make deploy-local   # Запустить всё
make docker-logs    # Проверить логи
```

### Изменение кода
```bash
# Код изменяется на хосте
# Контейнер автоматически перезапускается (если настроен hot reload)

# Или перезапустить вручную
make restart
```

### Работа с БД
```bash
# Подключиться к БД
make db-shell

# Создать бэкап
make backup-db

# Восстановить из бэкапа
make restore-db FILE=backups/backup_20260414_182600.sql
```

### Отладка
```bash
# Логи приложения
make docker-logs-app

# Логи БД
make docker-logs-db

# Shell в контейнере приложения
make shell

# Статус контейнеров
make status
```

## ⏳ Что осталось доделать

### Критично
1. **Подключить pgxpool в main.go** - инициализировать пул соединений
2. **Передать репозитории в handlers** - связать всё вместе
3. **Исправить систему шаблонов** - наследование base layout
4. **Создать первого админа** - через миграцию или CLI

### Важно
5. **Реализовать handlers с данными** - использовать репозитории
6. **Добавить bcrypt и sessions** - полная аутентификация
7. **CRUD операции** - инструкторы, слоты, бронирования

### Можно отложить
8. OpenWeatherMap API
9. VK Bot API
10. Загрузка фотографий

## 🎯 Следующий шаг

Обновить `cmd/server/main.go` для подключения к БД через pgxpool:

```go
import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    // ...
)

func main() {
    cfg := config.Load()
    
    // Создать пул соединений
    pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()
    
    // Инициализировать репозитории
    bookingRepo := repository.NewBookingRepository(pool)
    // ...
}
```

Подробности в [NEXT_STEPS.md](NEXT_STEPS.md)
