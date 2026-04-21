# Posts Service

HTTP сервер на Go — сервис публикации постов.

## Стек

- **Go 1.21**
- **PostgreSQL** — хранение данных
- **JWT** — сессионные токены авторизации

## Архитектура

Три слоя с изоляцией через интерфейсы:

```
handler.go      — HTTP слой, приём запросов
service.go      — бизнес-логика
repository.go   — работа с БД
```

Все файлы в одном пакете `main`, без вложенных папок.

## Структура проекта

```
posts-service/
├── main.go
├── handler.go
├── service.go
├── repository.go
└── go.mod
```

## Запуск

```bash
go mod tidy

$env:ADDR=":8081"
$env:DSN="host=localhost port=5432 user=postgres password=ПАРОЛЬ dbname=posts sslmode=disable"

go run .
```

## Эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/test` | Проверка работы сервера |
| POST | `/register` | Регистрация пользователя |
| POST | `/login` | Авторизация, получение JWT токена |
| POST | `/dbtest` | Запись строки в БД |
| GET | `/messages` | Чтение всех записей из БД |

## Примеры запросов

**Регистрация**
```json
POST /register
{"username": "dasha", "password": "1234"}
```

**Авторизация**
```json
POST /login
{"username": "dasha", "password": "1234"}
```
Ответ — JWT токен:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## База данных

Таблицы создаются автоматически при старте сервера:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    text TEXT
);
```

## Выполненные лабораторные работы

**ЛР1** — HTTP сервер с чистой архитектурой и graceful shutdown  
**ЛР2** — Подключение к PostgreSQL, инициализация таблиц, операции чтения и записи  
**ЛР3** — Регистрация пользователей и авторизация с выдачей JWT токена