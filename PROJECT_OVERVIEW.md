# SUP-Anapa - Краткий обзор проекта

## ✅ Что готово

### Инфраструктура
- ✅ Полная структура Go проекта
- ✅ PostgreSQL схема с миграциями
- ✅ Docker + Docker Compose setup с изолированной сетью
- ✅ **Makefile для управления всей инфраструктурой (Docker-first)**
- ✅ Все инструменты работают внутри Docker
- ✅ pgxpool для подключения к БД
- ✅ .gitignore, .env.example

### Backend
- ✅ Модели данных (Instructor, Slot, Booking, Admin, WeatherCache)
- ✅ Репозитории на pgxpool с context.Context
- ✅ Сервисы (BookingService, WeatherService, NotificationService)
- ✅ HTTP handlers (базовые заглушки)
- ✅ Middleware для аутентификации
- ✅ Chi router настроен

### Frontend
- ✅ HTML шаблоны для всех страниц:
  - Главная страница с галереей
  - Страница бронирования
  - Админка (dashboard, инструкторы, слоты, бронирования)
  - Страница входа
- ✅ Tailwind CSS подключен
- ✅ htmx подключен для динамики

### Документация
- ✅ Техническое задание (TECHNICAL_SPEC.md)
- ✅ Документация по разработке (DEVELOPMENT.md)
- ✅ README с инструкциями
- ✅ Руководство по Makefile (MAKEFILE_GUIDE.md)
- ✅ Docker инфраструктура (DOCKER_INFRASTRUCTURE.md)
- ✅ Следующие шаги (NEXT_STEPS.md)

## 🚀 Быстрый старт

```bash
# 1. Создать .env файл
make setup

# 2. Запустить всё одной командой (Docker + БД + миграции)
make deploy-local

# Приложение доступно на http://localhost:8080
```

**Важно:** Ничего не устанавливается на хост-машину! Всё работает в Docker.

## 📋 Основные команды Makefile

```bash
make help              # Показать все команды
make deploy-local      # Локальный деплой (Docker + миграции)
make deploy-prod       # Продакшн деплой
make docker-logs       # Логи всех контейнеров
make docker-logs-app   # Логи приложения
make migrate-up        # Применить миграции (внутри контейнера)
make db-shell          # Подключиться к БД
make shell             # Shell в контейнере приложения
make backup-db         # Создать бэкап БД
```

## ⏳ Что нужно доделать для MVP

### Критично (без этого не запустится)
1. **Подключить БД к handlers** - добавить инициализацию репозиториев в main.go
2. **Исправить систему шаблонов** - сейчас не работает наследование base layout
3. **Создать первого админа** - через миграцию или CLI команду

### Важно (для базовой работы)
4. **Реализовать handlers с данными** - сейчас только заглушки
5. **Доделать аутентификацию** - добавить bcrypt и gorilla/sessions
6. **CRUD для инструкторов и слотов** - подключить к репозиториям

### Можно отложить
7. Интеграция OpenWeatherMap API (пока моковые данные)
8. VK Bot API (пока логирование)
9. Загрузка фотографий

Подробности в [NEXT_STEPS.md](NEXT_STEPS.md)

## 📁 Структура проекта

```
sup-anapa/
├── cmd/server/              # Точка входа
├── internal/
│   ├── config/             # Конфигурация
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Auth middleware
│   ├── models/             # Модели данных
│   ├── repository/         # Работа с БД
│   └── services/           # Бизнес-логика
├── web/
│   ├── templates/          # HTML шаблоны
│   │   ├── layouts/       # Base layout
│   │   ├── public/        # Публичные страницы
│   │   └── admin/         # Админка
│   └── static/            # CSS, JS, изображения
├── migrations/             # SQL миграции
├── Makefile               # Управление инфраструктурой
├── docker-compose.yml     # Docker setup
└── Dockerfile             # Образ приложения
```

## 🗄️ База данных

### Таблицы
- **instructors** - Инструкторы (имя, фото, описание, телефон)
- **slots** - Временные слоты (дата, время, цена, инструктор)
- **bookings** - Бронирования (клиент, слот, статус)
- **admins** - Администраторы (логин, пароль)
- **weather_cache** - Кэш погоды

### Миграции
```bash
make migrate-up      # Применить
make migrate-down    # Откатить
make migrate-status  # Статус
```

## 🔧 Технологии

- **Backend**: Go 1.21+, Chi router
- **Database**: PostgreSQL 15+
- **Frontend**: Go templates, htmx, Tailwind CSS
- **Tools**: Docker, Make, golang-migrate, air
- **APIs**: OpenWeatherMap, VK Bot API (planned)

## 📝 Переменные окружения

Создайте `.env` файл:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/sup_anapa?sslmode=disable
WEATHER_API_KEY=your_openweathermap_api_key
VK_BOT_TOKEN=your_vk_bot_token
SESSION_SECRET=random_secret_string
```

## 🎯 Roadmap

### v0.1 - MVP
- [ ] Подключение БД к handlers
- [ ] Аутентификация админки
- [ ] CRUD инструкторов и слотов
- [ ] Создание бронирований

### v0.2 - Погода
- [ ] Интеграция OpenWeatherMap
- [ ] Кэширование погоды
- [ ] Отображение при бронировании

### v0.3 - Уведомления
- [ ] VK Bot API интеграция
- [ ] Подтверждение броней
- [ ] Статистика в админке

### v1.0 - Production Ready
- [ ] Загрузка фотографий
- [ ] Email уведомления
- [ ] Тесты
- [ ] Оптимизация

## 📚 Документация

- [README.md](README.md) - Основная информация
- [TECHNICAL_SPEC.md](TECHNICAL_SPEC.md) - Техническое задание
- [DEVELOPMENT.md](DEVELOPMENT.md) - Документация по разработке
- [MAKEFILE_GUIDE.md](MAKEFILE_GUIDE.md) - Руководство по Makefile
- [NEXT_STEPS.md](NEXT_STEPS.md) - Следующие шаги

## 🎉 Готово к работе!

Проект полностью настроен и готов к разработке. Используйте Makefile для управления всей инфраструктурой - от локальной разработки до продакшн деплоя.

```bash
make help  # Начните отсюда!
```
