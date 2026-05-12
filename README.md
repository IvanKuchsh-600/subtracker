# SubTracker

REST сервис для управления подписками пользователей.

## Запуск

```bash
docker-compose up --build

Сервис будет доступен: http://localhost:8080

API
Метод	URL	Описание
POST	/api/v1/subscriptions	Создать подписку
GET	/api/v1/subscriptions	Получить все подписки
GET	/api/v1/subscriptions/{id}	Получить подписку по ID
PUT	/api/v1/subscriptions/{id}	Обновить подписку
DELETE	/api/v1/subscriptions/{id}	Удалить подписку
GET	/api/v1/subscriptions/total-cost	Подсчитать стоимость за период

Примеры запросов

Создать подписку
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'

Подсчитать стоимость за 2025 год
curl -X GET "http://localhost:8080/api/v1/subscriptions/total-cost?from_date=01-2025&to_date=12-2025"

Обновить цену подписки
curl -X PUT "http://localhost:8080/api/v1/subscriptions/{id}" \
  -H "Content-Type: application/json" \
  -d '{"price": 500}'

Удалить подписку
curl -X DELETE "http://localhost:8080/api/v1/subscriptions/{id}"

Swagger документация
Открой в браузере: http://localhost:8080/swagger/index.html

Технологии
Go 1.25.5

Gin

PostgreSQL 16

pgxpool

Docker / Docker Compose

Swagger


Структура проекта
subtracker/
├── cmd/          # Точка входа
├── internal/
│   ├── domain/          # Бизнес-сущности
│   ├── http/            # Handlers, DTO, Router
│   ├── repository/      # Работа с БД
│   ├── usecase/         # Бизнес-логика
│   └── config/          # Конфигурация
├── migrations/          # SQL миграции
├── docs/                # Swagger
├── docker-compose.yml
├── Dockerfile
└── config.yaml