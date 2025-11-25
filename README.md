# Short URL Service

Сервис для сокращения URL на Go с поддержкой PostgreSQL и Docker

## Возможности

- **Сокращение URL** - преобразование длинных ссылок в короткие алиасы
- **Редирект** - автоматическое перенаправление по коротким ссылкам
- **REST API** - полный набор endpoints для управления URL
- **Basic Auth** - защита API с помощью аутентификации
- **Поддержка кастомных алиасов** - возможность задавать свои короткие коды
- **Контейнеризация** - полная поддержка Docker и Docker Compose
- **Тестирование** - комплексные unit-тесты с моками

## Технологии

- **Go 1.25.3+** - язык разработки
- **PostgreSQL** - основное хранилище данных
- **Chi Router** - HTTP роутер
- **Docker & Docker Compose** - контейнеризация и оркестрация
- **Testify & Mockery** - тестирование и генерация моков
- **Cleanenv** - управление конфигурацией

## Быстрый старт

### 1. Клонирование репозитория
```
bash
git clone https://github.com/Alex322322/short-url.git
cd short-url
```

### 2. Запуск с Docker Compose
```
bash
docker compose up -d --build
```

### 3. Использование API
Создание короткой ссылки
```
bash
curl -X POST http://localhost:8084/url \
  -u username:userpass \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/very/long/url/path",
    "alias": "my-link"
  }'
```
Ответ:
```
json
{
  "status": "OK",
  "alias": "my-link"
}
```
Редирект по короткой ссылке
```
bash
curl -L http://localhost:8084/my-link
```
Или просто открой в браузере: http://localhost:8084/my-link

Удаление ссылки
```
bash
curl -X DELETE http://localhost:8084/url/my-link \
  -u username:userpass
```

## API Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| `POST` | `/url` | Создание короткой ссылки | Basic Auth |
| `GET` | `/{alias}` | Редирект по короткой ссылке | Public |
| `DELETE` | `/url/{alias}` | Удаление ссылки | Basic Auth |

### Примеры запросов
Создание ссылки с автоматическим алиасом
```
bash
curl -X POST http://localhost:8084/url \
  -u username:userpass \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```

Создание ссылки с кастомным алиасом
```
bash
curl -X POST http://localhost:8084/url \
  -u username:userpass \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/Alex322322/short-url",
    "alias": "github-project"
  }'
```

## Разработка
Локальная разработка
Установи зависимости:
```
bash
go mod download
```
Запусти PostgreSQL:
```
bash
docker compose up pg -d
```
Настрой окружение:
```
bash
cp .env.example .env
# отредактируй .env файл
```
Запусти приложение:
```
bash
go run ./cmd/app
```
Запуск тестов
```
bash
# Все тесты
go test ./...

# Тесты с покрытием
go test -cover ./...
```
Генерация моков
```
bash
go generate ./...
```

## Структура проекта
| Директория | Назначение |
|------------|------------|
| `cmd/app/` | Главный исполняемый файл приложения |
| `internal/config/` | Конфигурация и переменные окружения |
| `internal/http/server/` | HTTP handlers, middleware, роутинг |
| `internal/lib/api/` | Вспомогательные функции для API |
| `internal/lib/logger/` | Настройка логирования |
| `internal/lib/random/` | Генерация случайных алиасов |
| `internal/storage/` | Интерфейсы хранилища данных |
| `storage/postgres/` | Реализация для PostgreSQL |
| `storage/sqlite/` | Реализация для SQLite |
| `tests/` | Функциональные тесты |

## Безопасность
API защищено Basic Authentication

Пароли хранятся в переменных окружения

Валидация входных данных

Подготовленные SQL statements для защиты от инъекций
