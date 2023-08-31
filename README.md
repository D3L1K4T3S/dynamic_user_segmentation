# Dynamic user segmentation

Микросервис для работы с сегментами пользователей. Возможность получать сегменты c ttl, к которым относится пользователь.
Также микросервис хранит всю историю действий совершенную над сегментами

## Используемые технологии

1. PostgreSQL - как хранилище данных. Драйвер (pgx)
2. Swagger - для документации API
3. Postman - для создания HTTP запросов
3. ECHO - фреймворк для развертывания контроллера
4. Docker - для создания контейнера
5. Docker-compose - инструмент для оркестрации контейнеров
6. golang-migrate/migrate - для миграции базы данных
7. golang/mock - создание моков для тестирования сервисов
8. pgxmock - создания моков для тестирования репозитория

``
Был реализован Graceful Shutdown для завершения работы сервиса
``

## Complete
- Добавление/удаление сегментов
- Получение всех сегментов хранящиеся в базе данных
- Автоматическое добавление сегментов пользователю при указании сегменту определенного %
- Добавление/удаление сегментов у пользователя
- Получение всех сегментов конкретного пользователя
- Обновление TTL для сегмента пользователя
- Реализовано актуальное получение сегментов
- Хранение истории операций, которые были совершены над сегментами пользователя

### Diagram database
База данных была нормализована и нет связей many to many


## Decisions


## Getting Started

Для запуска необходимо:
- Настроить .env файл в директорию с проектом. Заполнить как также, как представлен в этом проекте
- Настроить конфигурацию работы сервера в config/config.yaml

## Usage

- Запустить сервер: `make compose-up`
- Для запусков тестов: `make test`
- Для запусков тестов с покрытием: `make cover-html`
- Для запуска линтера: `make linter-golangci`
- Документация можно будет посмотреть `http://{host}:{port}/swagger/index.html`

## Examples
