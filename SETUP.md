# Инструкция по настройке SUP Anapa

## Требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

## Быстрый старт

1. Клонируйте репозиторий
2. Создайте файл `.env` в корне проекта:

```env
WEATHER_API_KEY=your_weather_api_key
VK_BOT_TOKEN=your_vk_bot_token
SESSION_SECRET=your_random_secret_key_min_32_chars
```

3. Запустите приложение:

```bash
docker compose up -d
```

4. Приложение будет доступно по адресу: http://localhost:8080

## Доступ к админке

После первого запуска автоматически создается администратор:

- **URL**: http://localhost:8080/admin/login
- **Логин**: `admin`
- **Пароль**: `admin123`

⚠️ **Важно**: Измените пароль администратора после первого входа!

## Структура проекта

```
sup-anapa/
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── config/         # Конфигурация
│   ├── handlers/       # HTTP обработчики
│   ├── models/         # Модели данных
│   ├── repository/     # Работа с БД
│   └── services/       # Бизнес-логика
├── migrations/         # SQL миграции
├── web/
│   ├── static/        # Статические файлы
│   └── templates/     # HTML шаблоны
└── docker-compose.yml
```

## База данных

Приложение использует PostgreSQL. Миграции применяются автоматически при запуске.

### Подключение к БД

```bash
docker compose exec postgres psql -U postgres -d sup_anapa
```

## Разработка

### Локальный запуск

```bash
# Установить зависимости
go mod download

# Запустить БД
docker compose up -d postgres

# Запустить приложение
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/sup_anapa?sslmode=disable"
export PORT=8080
go run cmd/server/main.go
```

### Создание новой миграции

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

## Полезные команды

```bash
# Просмотр логов
docker compose logs -f app

# Перезапуск приложения
docker compose restart app

# Остановка всех сервисов
docker compose down

# Пересборка образа
docker compose build --no-cache
```

## Функционал

### Публичная часть
- Главная страница с информацией о SUP прогулках
- Страница бронирования с выбором инструктора и времени
- Интеграция с прогнозом погоды

### Админ-панель
- Управление инструкторами
- Управление временными слотами
- Просмотр и управление бронированиями
- Статистика

## Технологии

- **Backend**: Go, Chi Router
- **Database**: PostgreSQL, pgx
- **Frontend**: HTML, Tailwind CSS, HTMX
- **Auth**: bcrypt, gorilla/sessions
- **Deployment**: Docker, Docker Compose
