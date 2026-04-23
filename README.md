# Posts Service

HTTP сервер на Go — сервис публикации постов.

## Стек

- **Go 1.21**
- **PostgreSQL** — хранение данных
- **JWT** — сессионные токены авторизации
- **Docker Compose + Nginx** — контейнерный запуск API, БД и Web UI

## Архитектура

Три слоя с изоляцией через интерфейсы:

```
handler.go      — HTTP слой, приём запросов
middleware.go   — проверка JWT токена
service.go      — бизнес-логика
repository.go   — работа с БД
```

## Структура проекта

```
posts-service/
├── main.go
├── handler.go
├── middleware.go
├── service.go
├── repository.go
├── go.mod
├── docker-compose.yml
├── deploy/
│   └── nginx/
│       └── default.conf
└── web/
    ├── index.html
    ├── styles.css
    └── app.js
```

## Запуск через Docker Compose (рекомендуется)

```bash
docker compose up --build
```

После запуска:

- **Web UI:** `http://localhost:8080`
- **Go API:** `http://localhost:8081`
- **PostgreSQL:** `localhost:5432`

Остановка:

```bash
docker compose down
```

## Локальный запуск без Docker

1. Убедитесь, что PostgreSQL запущен локально и доступен по DSN.
2. Запустите сервер:

```powershell
go mod tidy

$env:ADDR=":8081"
$env:DSN="host=localhost port=5432 user=postgres password=ПАРОЛЬ dbname=posts sslmode=disable"

go run .
```

3. Откройте `web/index.html` в браузере.
   - Если открываете как `file://`, интерфейс отправляет запросы в `http://localhost:8081`.
   - Если запуск через Docker Compose, интерфейс доступен на `http://localhost:8080` и работает через nginx-прокси.

## Проверка через Web UI

Интерфейс позволяет протестировать все endpoint:

1. **Регистрация** → `POST /register`
2. **Логин** → `POST /login` (JWT сохраняется автоматически)
3. **Создать пост** → `POST /posts` (с `Authorization: Bearer <token>`)
4. **Показать мои посты** → `GET /posts` (возвращает посты текущего пользователя)

## Эндпоинты

| Метод | URL | Защита | Описание |
|-------|-----|--------|----------|
| GET | `/test` | — | Проверка работы сервера |
| POST | `/register` | — | Регистрация пользователя |
| POST | `/login` | — | Авторизация, получение JWT токена |
| POST | `/posts` | JWT | Создание поста |
| GET | `/posts` | JWT | Просмотр своих постов |

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

**Создание поста** (требует токен)
```
POST /posts
Authorization: Bearer <токен>

{"content": "мой первый пост"}
```

**Просмотр своих постов** (требует токен)
```
GET /posts
Authorization: Bearer <токен>
```

## База данных

Таблицы создаются автоматически при старте сервера:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL
);
```

## Выполненные лабораторные работы

**ЛР1** — HTTP сервер с чистой архитектурой и graceful shutdown  
**ЛР2** — Подключение к PostgreSQL, инициализация таблиц, операции чтения и записи  
**ЛР3** — Регистрация пользователей и авторизация с выдачей JWT токена  
**ЛР4** — Хэндлеры создания и просмотра постов  
**ЛР5** — Middleware для проверки JWT и передачи user_id в хэндлеры