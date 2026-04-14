# Следующие шаги для завершения MVP

## Что уже готово ✅
- Полная структура проекта
- Модели данных и SQL миграции
- Репозитории для работы с БД
- Сервисы (booking, weather, notification)
- HTML шаблоны для всех страниц
- Docker setup
- Базовые handlers (заглушки)

## Что нужно доделать для рабочего MVP

### 1. Подключение базы данных к handlers
**Файл**: `cmd/server/main.go`

Нужно:
- Подключиться к PostgreSQL
- Инициализировать репозитории
- Передать их в handlers

```go
// Пример
db, err := sql.Open("postgres", cfg.DatabaseURL)
bookingRepo := repository.NewBookingRepository(db)
// и т.д.
```

### 2. Реализовать handlers с реальной логикой
**Файлы**: `internal/handlers/*.go`

Сейчас handlers возвращают только шаблоны. Нужно:
- Получать данные из БД
- Передавать их в шаблоны
- Обрабатывать POST запросы

### 3. Создать первого администратора
**Способ 1**: Через миграцию
```sql
-- migrations/002_create_admin.up.sql
INSERT INTO admins (username, password_hash) 
VALUES ('admin', '$2a$10$...');  -- bcrypt hash пароля
```

**Способ 2**: Через CLI команду
```go
// cmd/createadmin/main.go
```

### 4. Реализовать аутентификацию
**Нужно**:
- Установить `golang.org/x/crypto/bcrypt`
- Установить `github.com/gorilla/sessions`
- Доработать `internal/middleware/auth.go`
- Доработать `handlers.AdminLoginPost`

### 5. Исправить шаблоны
**Проблема**: Сейчас шаблоны используют `{{define}}`, но не наследуют base layout

**Решение**: Переделать систему шаблонов:
```go
// Вариант 1: Использовать template.ParseFiles с несколькими файлами
// Вариант 2: Использовать {{template "base" .}} внутри шаблонов
```

### 6. Добавить HTMX endpoints
Для динамических обновлений нужны API endpoints:
- `GET /api/slots?from=&to=` - получить слоты
- `POST /api/instructors` - создать инструктора
- `PATCH /api/bookings/:id/status` - обновить статус

## Быстрый старт для тестирования

```bash
# 1. Запустить PostgreSQL
docker-compose up -d postgres

# 2. Применить миграции
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/sup_anapa?sslmode=disable" up

# 3. Создать .env
cp .env.example .env

# 4. Запустить приложение
go run cmd/server/main.go
```

## Приоритет задач

### Критично (без этого не запустится)
1. Подключение БД в main.go
2. Исправление системы шаблонов
3. Создание первого админа

### Важно (для базовой работы)
4. Реализация handlers с данными из БД
5. Аутентификация админки
6. CRUD для инструкторов и слотов

### Можно отложить
7. Интеграция OpenWeatherMap (пока моковые данные)
8. VK Bot API (пока логирование)
9. Загрузка фотографий

## Полезные команды

```bash
# Проверка компиляции
go build ./...

# Запуск
go run cmd/server/main.go

# Тесты
go test ./...

# Форматирование
go fmt ./...
```

## Контакты

Если нужна помощь с реализацией - пиши!
