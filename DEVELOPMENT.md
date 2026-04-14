# SUP Anapa - Документация по разработке

## Обзор проекта

SUP Anapa - веб-сервис для бронирования SUP-прогулок в Анапе. Проект построен на Go с использованием PostgreSQL и минималистичного фронтенда на Go templates + htmx.

## Текущий статус

### Реализовано
- ✅ Структура проекта
- ✅ Модели данных (instructors, slots, bookings, admins, weather_cache)
- ✅ SQL миграции для PostgreSQL
- ✅ Базовые репозитории для работы с БД
- ✅ Сервисы (booking, weather, notification)
- ✅ HTTP handlers (заглушки)
- ✅ HTML шаблоны для всех страниц
- ✅ Docker setup (Dockerfile, docker-compose.yml)
- ✅ Middleware для аутентификации (базовая версия)

### В процессе
- 🔄 Интеграция с базой данных
- 🔄 Полная реализация handlers

### Планируется
- ⏳ Интеграция с OpenWeatherMap API
- ⏳ Интеграция с VK Bot API
- ⏳ Система сессий для админки
- ⏳ Хеширование паролей (bcrypt)
- ⏳ HTMX интеграция для динамических обновлений

## Архитектура

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│      Chi Router (main.go)       │
│  ┌───────────┬────────────────┐ │
│  │  Public   │     Admin      │ │
│  │  Routes   │    Routes      │ │
│  └─────┬─────┴────────┬───────┘ │
└────────┼──────────────┼─────────┘
         │              │
         ▼              ▼
┌─────────────┐  ┌──────────────┐
│  Handlers   │  │  Middleware  │
└──────┬──────┘  └──────────────┘
       │
       ▼
┌─────────────────────────────────┐
│          Services               │
│  ┌──────────┬─────────────────┐ │
│  │ Booking  │ Weather │ Notif │ │
│  └────┬─────┴─────────────────┘ │
└───────┼─────────────────────────┘
        │
        ▼
┌─────────────────────────────────┐
│        Repositories             │
│  ┌──────────┬──────────────────┐│
│  │ Booking  │ Slot │ Instructor││
│  └────┬─────┴──────────────────┘│
└───────┼─────────────────────────┘
        │
        ▼
┌─────────────────────────────────┐
│         PostgreSQL              │
└─────────────────────────────────┘
```

## Структура базы данных

### Таблицы

1. **instructors** - Инструкторы
   - id, name, photo, description, phone
   - timestamps

2. **slots** - Временные слоты
   - id, date, start_time, end_time, price, max_people
   - instructor_id (FK)
   - timestamps

3. **bookings** - Бронирования
   - id, slot_id (FK), client_name, client_phone, client_email
   - people_count, status (pending/confirmed/cancelled)
   - timestamps

4. **admins** - Администраторы
   - id, username, password_hash
   - created_at

5. **weather_cache** - Кэш погоды
   - id, date, air_temp, water_temp, wind_speed, cloud_cover
   - description, cached_at

## Запуск проекта

### Быстрый старт с Makefile

```bash
# Установить инструменты
make install-tools

# Локальный деплой (всё в одной команде)
make deploy-local

# Запустить приложение
make dev  # с hot reload
# или
make run  # обычный запуск
```

### Продакшн деплой

```bash
make deploy-prod
```

### Все команды

Полный список команд:
```bash
make help
```

Подробная документация: [MAKEFILE_GUIDE.md](MAKEFILE_GUIDE.md)

## API Endpoints

### Публичные
- `GET /` - Главная страница
- `GET /booking` - Страница бронирования
- `POST /booking` - Создание бронирования

### Админка
- `GET /admin/login` - Страница входа
- `POST /admin/login` - Аутентификация
- `GET /admin` - Dashboard
- `GET /admin/instructors` - Управление инструкторами
- `POST /admin/instructors` - Создание инструктора
- `GET /admin/slots` - Управление слотами
- `POST /admin/slots` - Создание слота
- `GET /admin/bookings` - Просмотр броней
- `PATCH /admin/bookings/:id` - Обновление статуса брони

## Следующие шаги

### Приоритет 1 (Критично для MVP)
1. Подключить базу данных к handlers
2. Реализовать создание первого админа через CLI/миграцию
3. Реализовать систему сессий (gorilla/sessions)
4. Реализовать хеширование паролей (bcrypt)
5. Подключить HTMX для динамических обновлений

### Приоритет 2 (Основной функционал)
1. Интеграция OpenWeatherMap API
2. Кэширование погодных данных
3. Отображение доступных слотов на странице бронирования
4. CRUD операции для инструкторов и слотов

### Приоритет 3 (Дополнительно)
1. Интеграция VK Bot API
2. Email уведомления
3. Загрузка фотографий
4. Статистика в админке
5. Экспорт данных

## Технологии

- **Backend**: Go 1.21+, Chi router
- **Database**: PostgreSQL 15+, database/sql
- **Frontend**: Go html/template, htmx, Tailwind CSS
- **Deployment**: Docker, docker-compose
- **Migrations**: golang-migrate

## Переменные окружения

```env
PORT=8080
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
WEATHER_API_KEY=your_openweathermap_key
VK_BOT_TOKEN=your_vk_bot_token
SESSION_SECRET=random_secret_string
```

## Полезные команды

```bash
# Запуск тестов
go test ./...

# Форматирование кода
go fmt ./...

# Проверка кода
go vet ./...

# Сборка
go build -o bin/server cmd/server/main.go

# Миграции
migrate -path migrations -database $DATABASE_URL up
migrate -path migrations -database $DATABASE_URL down

# Docker
docker-compose up -d
docker-compose logs -f app
docker-compose down
```

## Контакты и поддержка

Для вопросов и предложений создавайте issues в репозитории проекта.
