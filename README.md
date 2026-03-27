# Тестовое задание для стажёра Backend

В компании нужен единый сервис бронирования переговорок: администраторы создают переговорки и настраивают расписание их доступности (например, по дням недели и времени), система сама формирует слоты для бронирования на основе расписания. Сотрудники просматривают свободные слоты и создают или отменяют брони.

Подробнее: [TASK.md](TASK.md)

## Технологический стек

| Категория | Технология |
|-----------|-----------|
| Язык | Go 1.26 |
| HTTP | net/http |
| БД | PostgreSQL 18 |
| Драйвер БД | pgx/v5 + pgxpool |
| Миграции | golang-migrate |
| Логирование | uber-zap |
| Контейнеризация | Docker (multi-stage build) + Docker Compose |
| Тесты | testify + mockery |
| Swagger | swaggo/swag |
| Конфигурация | envconfig |
| Валидация |go-playground/validator |

## Архитектура

Проект построен с разделением на слои:

```
HTTP Request
     │
 Middleware   ← auth, logging, recovery, trace
     │
 Handlers     ← преобразованние в DTO, возврат ответа
     │
 Service      ← бизнес-логика, валидация
     │
 Repository      ← репозитории
     │
 PostgreSQL
```

### Структура проекта

```
.
├── cmd
│   └── roomsbooking          # Точка входа
├── docs
├── internal
│   ├── core                 
│   │   ├── domain             # Доменные сущности(модели, ошибки)
│   │   ├── logger             # Инициализация логгера
│   │   ├── middleware         # auth, logging, recovery, trace
│   │   ├── repository
│   │   │   └── postgres       # Инициализация пула бд
│   │   └── transport
│   │       ├── response       # Отправка ответов хендлеров
│   │       └── server         # Настройка сервера
│   ├── features
│   │   ├── auth
│   │   │   ├── repository
│   │   │   │   └── postgres   # Репозиторий пользователей
│   │   │   ├── service        # Бизнес-логика авторизации
│   │   │   └── transport
│   │   │       └── http       # Хендлеры авторизации
│   │   ├── bookings
│   │   │   ├── respository
│   │   │   │   └── postgres   # Репозиторий бронирований
│   │   │   ├── service        # Бизнес-логика бронирований
│   │   │   └── transport
│   │   │       └── http       # Хендлеры бронирований
│   │   ├── rooms
│   │   │   ├── repository
│   │   │   │   └── postgres   # Репозиторий переговорок
│   │   │   ├── service        # Бизнес-логика переговорок
│   │   │   └── transport
│   │   │       └── http       # Хендлеры переговорок
│   │   ├── schedules
│   │   │   ├── repository
│   │   │   │   └── postgres   # Репозиторий расписания
│   │   │   ├── service
│   │   │   └── transport
│   │   │       └── http       # Хендлеры расписания
│   │   └── slots
│   │       ├── repository     
│   │       │   └── postgres   # Репозиторий слотов
│   │       ├── service        # Бизнес-логика слотов
│   │       └── transport
│   │           └── http       # Хендлеры слотов
│   └── utils
│       ├── hash               # Утилиты хеширования паролей
│       └── jwt_utils          # Утилиты генерации jwt токенов
├── migrations                 # Миграции
├── mocks                      # Моки для тестирования
│   ├── auth_service
│   ├── bookings_service
│   ├── schedules_service
│   └── slots_service
└── tests
    ├── e2e                    # e2e тесты
    │   
    └── unit                   # unit тесты
        ├── config
        ├── service
        │   ├── auth
        │   ├── bookings
        │   ├── schedules
        │   └── slots
        └── utils
```

## Запуск проекта

### Требования

- Docker и Docker Compose
- Go 1.26

### Запуск

```bash

  make up

```

### Сервисы

| Сервис | URL |
|--------|--------------|
| API | http://localhost:8080 |
| Swagger | http://localhost:80880/swagger |

## Тестирование

Код покрыт юнит тестами на 40.9%

```bash
# Unit-тесты 
make test-unit

# Unit-тесты с покрытием
make test-unit-coverage-html

# Unit-тесты с подсчетом покрытия
make test-unit-coverage-total

# E2E тесты (автоматически поднимает тестовое окружение)
make test-e2e

```


## Make-команды

```bash

make up              # Запуск проекта и всех зависимостей
make seed            # Наполнение проекта тестовыми данными

```

## Сделанные дополнительные задания

* **Регистрация и авторизация по email/паролю.** Реализовать эндпоинты `/register` и `/login` с выдачей JWT
* **Makefile.** Написать `Makefile` с командами: запуск проекта со всеми зависимостями (`make up`) и наполнение БД тестовыми данными (`make seed`).
* **Swagger-документация.** Настроить кодогенерацию Swagger-документации на основе аннотаций в коде.

[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)
