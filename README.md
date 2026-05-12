# SubTracker

REST сервис для управления подписками пользователей.

## Запуск

docker-compose up --build

Сервис запускается на http://localhost:8080

## API методы

| Метод | URL | Описание |
|-------|-----|----------|
| POST | /api/v1/subscriptions | Создать подписку |
| GET | /api/v1/subscriptions | Получить все подписки |
| GET | /api/v1/subscriptions/{id} | Получить подписку по ID |
| PUT | /api/v1/subscriptions/{id} | Обновить подписку |
| DELETE | /api/v1/subscriptions/{id} | Удалить подписку |
| GET | /api/v1/subscriptions/total-cost | Сумма подписок за период |

## Пример запроса

curl -X POST http://localhost:8080/api/v1/subscriptions -H "Content-Type: application/json" -d '{"service_name": "Yandex Plus", "price": 400, "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba", "start_date": "07-2025"}'

## Документация

http://localhost:8080/swagger/index.html

## Технологии

Go 1.25.5, Gin, PostgreSQL, pgxpool, Docker, Swagger