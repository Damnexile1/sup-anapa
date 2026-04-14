# SUP-Anapa Makefile - Быстрая справка

## Первый запуск проекта

### 1. Установка инструментов
```bash
make install-tools
```
Установит golang-migrate и air для hot reload.

### 2. Локальный деплой (одна команда!)
```bash
make deploy-local
```
Эта команда:
- Проверит .env файл (создаст из .env.example если нужно)
- Установит зависимости Go
- Запустит PostgreSQL в Docker
- Применит миграции
- Соберет приложение

### 3. Запуск приложения
```bash
make run          # Обычный запуск
# или
make dev          # С hot reload (автоперезагрузка при изменениях)
```

## Продакшн деплой (одна команда!)

```bash
make deploy-prod
```
Эта команда:
- Проверит .env файл
- Остановит старые контейнеры
- Соберет Docker образ
- Запустит все сервисы
- Применит миграции

## Основные команды

### Разработка
```bash
make help           # Показать все доступные команды
make setup          # Первоначальная настройка проекта
make start          # Быстрый старт (setup + deploy + run)
make dev            # Запуск с hot reload
make build          # Собрать приложение
make test           # Запустить тесты
make fmt            # Форматировать код
make vet            # Проверить код
```

### Docker
```bash
make docker-up      # Запустить все сервисы
make docker-down    # Остановить все сервисы
make docker-logs    # Показать логи
make docker-restart # Перезапустить сервисы
make docker-clean   # Удалить контейнеры и volumes
```

### База данных
```bash
make db-up          # Запустить только PostgreSQL
make db-down        # Остановить PostgreSQL
make db-shell       # Подключиться к БД через psql
make backup-db      # Создать бэкап БД
make restore-db FILE=backups/backup.sql  # Восстановить БД
```

### Миграции
```bash
make migrate-up     # Применить все миграции
make migrate-down   # Откатить последнюю миграцию
make migrate-reset  # Сбросить и применить заново
make migrate-status # Показать статус миграций
make migrate-create NAME=add_users  # Создать новую миграцию
```

### Утилиты
```bash
make status         # Показать статус сервисов
make clean          # Очистить сборочные файлы
make restart        # Перезапустить приложение
```

## Типичные сценарии

### Начать работу над проектом
```bash
make setup
make deploy-local
make dev
```

### Деплой на сервер
```bash
# На сервере
git pull
make deploy-prod
make docker-logs  # Проверить логи
```

### Создать новую миграцию
```bash
make migrate-create NAME=add_reviews_table
# Отредактировать файлы в migrations/
make migrate-up
```

### Откатить изменения в БД
```bash
make migrate-down   # Откатить последнюю
# или
make migrate-reset  # Пересоздать всё
```

### Бэкап перед важными изменениями
```bash
make backup-db
# Делаем изменения...
# Если что-то пошло не так:
make restore-db FILE=backups/backup_20260414_181756.sql
```

## Переменные окружения

Создайте `.env` файл (будет создан автоматически из `.env.example`):

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/sup_anapa?sslmode=disable
WEATHER_API_KEY=your_openweathermap_api_key
VK_BOT_TOKEN=your_vk_bot_token
SESSION_SECRET=random_secret_string
```

## Порты

- **8080** - Веб-приложение
- **5432** - PostgreSQL

## Troubleshooting

### Ошибка "migrate not found"
```bash
make install-tools
```

### Ошибка подключения к БД
```bash
make db-up
# Подождать 3-5 секунд
make migrate-up
```

### Порт 8080 занят
```bash
# Изменить PORT в .env файле
# или остановить другое приложение
lsof -ti:8080 | xargs kill
```

### Очистить всё и начать заново
```bash
make docker-clean
make clean
make deploy-local
```

## Полезные алиасы для .bashrc/.zshrc

```bash
alias sup-start="cd ~/GolandProjects/sup-anapa && make start"
alias sup-dev="cd ~/GolandProjects/sup-anapa && make dev"
alias sup-logs="cd ~/GolandProjects/sup-anapa && make docker-logs"
alias sup-deploy="cd ~/GolandProjects/sup-anapa && make deploy-prod"
```
