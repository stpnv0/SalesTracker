# SalesTracker

Сервис для учёта финансовых операций (доходы / расходы) с CRUD, SQL-аналитикой, CSV-экспортом и веб-интерфейсом.

---

## Возможности

- **CRUD операции** с финансовыми записями (транзакциями)
- **Аналитика** — расчет суммы, среднего, медианы, 90-го перцентиля
- **Группировка** по дням, неделям, месяцам и категориям
- **Фильтрация и сортировка** записей
- **Экспорт данных** в CSV
- **Веб-интерфейс** для управления записями

## Стек технологий

- **Go 1.25**
- **wbf** — внутренний фреймворк (ginext, dbpg, config, zlog, helpers, retry)
- **PostgreSQL 17** — основное хранилище
- **goose** — миграции
- **go-playground/validator** — валидация DTO
- **Docker / Docker Compose** — контейнеризация


---

## Структура проекта

```
.
├── cmd/server/           # Точка входа
├── internal/
│   ├── app/              # Сборка и жизненный цикл приложения
│   ├── config/           # Конфигурация
│   ├── domain/           # Доменные модели и ошибки
│   ├── handler/          # HTTP handlers
│   ├── service/          # Бизнес-логика
│   ├── repository/       # Работа с PostgreSQL
│   ├── router/           # Маршрутизация
│   ├── middleware/       # CORS, Logging, RequestID
│   └── export/           # Экспорт CSV
├── web/                  # Веб-интерфейс (HTML, CSS, JS)
├── migrations/           # SQL миграции
├── Dockerfile            # Контейнеризация
└── docker-compose.yaml   # Запуск с Postgres
```

---
## API

### CRUD

| Метод    | Путь             | Описание           |
|----------|------------------|--------------------|
| `POST`   | `/api/items`     | Создать запись     |
| `GET`    | `/api/items`     | Список с фильтрами |
| `GET`    | `/api/items/:id` | Получить по ID     |
| `PUT`    | `/api/items/:id` | Обновить запись    |
| `DELETE` | `/api/items/:id` | Удалить запись     |

#### Query-параметры для GET /api/items

| Параметр   | Тип                            | Описание               |
|------------|--------------------------------|------------------------|
| `from`     | `YYYY-MM-DD`                   | Начальная дата         |
| `to`       | `YYYY-MM-DD`                   | Конечная дата          |
| `category` | `string`                       | Фильтр по категории    |
| `type`     | `income\|expense`              | Фильтр по типу         |
| `sort_by`  | `date\|amount\|category\|type` | Поле сортировки        |
| `order`    | `asc\|desc`                    | Направление сортировки |
| `limit`    | `int`                          | Лимит записей          |
| `offset`   | `int`                          | Смещение               |


### Аналитика

| Метод   | Путь                                                   | Описание                  |
|---------|--------------------------------------------------------|---------------------------|
| `GET`   | `/api/analytics?from=...&to=...&group_by=...&type=...` |  Агрегированная аналитика |

#### Query-параметры

| Параметр   | Обязательный | Описание                           |
|------------|--------------|------------------------------------|
| `from`     | да           | Начало периода (`YYYY-MM-DD`)      |
| `to`       | да           | Конец периода (`YYYY-MM-DD`)       |
| `group_by` | нет          | `day`, `week`, `month`, `category` |
| `type`     | да           | Тип операции (`income`/`expense`)      |

### Экспорт

| Метод   | Путь                                | Описание               |
|---------|-------------------------------------|------------------------|
| `GET`   | `/api/export/csv?from=...&to=...`   | Скачать данные в CSV   |

Поддерживает те же фильтры: `from`, `to`, `category`, `type`.

---

## Запуск

## Запуск

### Docker Compose

```bash
# Клонировать
git clone https://github.com/stpnv0/SalesTracker.git
cd EventBooker

# Запустить
docker-compose up --build

# Открыть
http://localhost:8080
```

## Структура БД

### Таблица `items`

| Колонка       | Тип             | Ограничения                                         |
|---------------|-----------------|-----------------------------------------------------|
| `id`          | `UUID`          | `PRIMARY KEY`                                       |
| `type`        | `VARCHAR(10)`   | `NOT NULL`, `CHECK (type IN ('income', 'expense'))` |
| `amount`      | `NUMERIC(15,2)` | `NOT NULL`, `CHECK (amount > 0)`                    |
| `category`    | `VARCHAR(100)`  | `NOT NULL`                                          |
| `description` | `TEXT`          | `NOT NULL DEFAULT ''`                               |
| `date`        | `DATE`          | `NOT NULL`                                          |
| `created_at`  | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                            |
| `updated_at`  | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                            |

