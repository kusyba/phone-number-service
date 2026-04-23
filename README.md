<div align="center">

# 📞 Phone Number Management Service

[![Go Version](https://img.shields.io/badge/Go-1.25.8-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=flat-square&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-FFB000?style=flat-square)](LICENSE)

**REST API сервис для управления телефонными номерами**  
Группы пользователей · Битовые права доступа · Защита от дубликатов

</div>

## 🌟 Возможности

### 📞 Управление номерами

| Функция | Описание |
|---------|----------|
| **E.164 Нормализация** | Приведение номеров к единому стандарту |
| **Массовый импорт** | Загрузка с защитой от дубликатов |
| **Гибкий поиск** | Поиск по номеру, стране, региону, провайдеру |
| **Фронтенд формат** | `+7 (916) 123-45-67` |
| **Пагинация** | `limit` / `offset` |

### 👥 Управление группами

| Операция | Метод | Эндпоинт |
|----------|-------|----------|
| Создание | `POST` | `/api/groups` |
| Список | `GET` | `/api/groups` |
| Получение | `GET` | `/api/groups/:id` |
| Обновление | `PUT` | `/api/groups/:id` |
| Удаление | `DELETE` | `/api/groups/:id` |
| Добавить участника | `POST` | `/api/groups/:id/users/:userId` |
| Удалить участника | `DELETE` | `/api/groups/:id/users/:userId` |
| Группы пользователя | `GET` | `/api/users/:userId/groups` |

### 🔐 Система прав

| Значение | Флаг | Описание |
|:--------:|------|----------|
| 16 | `can_change_phone` | Запрет на изменение номера |
| 32 | `can_change_email` | Запрет на изменение email |
| 64 | `can_change_tariff` | Запрет на смену тарифа |
| 128 | `can_change_relative_settings` | Запрет на изменение настроек родственника |
| 256 | `can_leave_corporation` | Запрет на исключение из корпорации |
| 512 | `can_disable_incident_alerts` | Запрет на отключение уведомлений |
| 1024 | `can_delete_account` | Запрет на удаление аккаунта |

> 💡 **Комбинирование:** `2047` = все права, `16+32` = два права

## 🛠 Технологии

| Категория | Технология |
|-----------|------------|
| **Язык** | Go 1.25.8 |
| **Web фреймворк** | Echo v4.12.0 |
| **База данных** | PostgreSQL 15 |
| **Миграции** | Goose |
| **Генерация запросов** | sqlc |
| **Контейнеризация** | Docker |
| **Логирование** | JSON |

## 🚀 Быстрый старт

### 🐳 Docker

```bash
git clone https://github.com/kusyba/phone-number-service.git
cd phone-number-service
docker-compose up --build
💻 Локально
bash
git clone https://github.com/kusyba/phone-number-service.git
cd phone-number-service
make install
make migrate-up
make generate
make run
📚 API Эндпоинты
<details> <summary><b>Номера телефонов</b></summary>
POST /api/numbers/import

bash
curl -X POST http://localhost:8080/api/numbers/import \
  -H "Content-Type: application/json" \
  -d '{"numbers":["+79161234567","89161234567","9161234567","123"],"source":"telegram"}'
json
{"accepted":1,"skipped":2,"errors":1}
GET /api/numbers/search

bash
curl "http://localhost:8080/api/numbers/search?number=916&limit=10&offset=0"
POST /api/phones/format

bash
curl -X POST http://localhost:8080/api/phones/format \
  -H "Content-Type: application/json" \
  -d '{"number":"89161234567"}'
json
{"formatted":"+7 (916) 123-45-67","original":"89161234567"}
GET /api/phones/:id/format

bash
curl http://localhost:8080/api/phones/1/format
json
{"formatted":"+7 (916) 123-45-67","id":1}
</details><details> <summary><b>Группы пользователей</b></summary>
POST /api/groups

bash
curl -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{"name":"admins","description":"Full access","flags":2047}'
GET /api/groups

bash
curl "http://localhost:8080/api/groups?limit=10&offset=0"
PUT /api/groups/:id

bash
curl -X PUT http://localhost:8080/api/groups/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"super_admins","flags":2047}'
DELETE /api/groups/:id

bash
curl -X DELETE http://localhost:8080/api/groups/1
POST /api/groups/:id/users/:userId

bash
curl -X POST http://localhost:8080/api/groups/1/users/100
DELETE /api/groups/:id/users/:userId

bash
curl -X DELETE http://localhost:8080/api/groups/1/users/100
GET /api/users/:userId/groups

bash
curl "http://localhost:8080/api/users/100/groups"
</details><details> <summary><b>Пользователь</b></summary>
GET /api/me

bash
curl -H "X-User-ID: 100" http://localhost:8080/api/me
json
{
  "user_id": 100,
  "groups": [...],
  "flags": 2047,
  "permissions": {...}
}
GET /health

bash
curl http://localhost:8080/health
json
{"status":"ok"}
</details>
🔒 Безопасность
Защита	Реализация
SQL-инъекции	Экранирование LIKE + параметризованные запросы
XSS	HTML-экранирование всех выходных данных
DDoS	Rate limiting (100 запросов/минуту)
Дубликаты	ON CONFLICT + транзакции
Graceful shutdown	Корректное завершение при SIGTERM
🧪 Тестирование
bash
# Unit-тесты нормализации
go test -v ./pkg/utils/

# Все тесты
go test -v ./...

# С покрытием
go test -cover ./...
📁 Структура проекта
text
phone-number-service/
├── cmd/server/           # Точка входа
├── internal/
│   ├── api/              # HTTP handlers
│   ├── config/           # Конфигурация
│   ├── database/         # Работа с БД + sqlc
│   ├── models/           # Модели данных
│   └── service/          # Бизнес-логика
├── pkg/
│   ├── logger/           # Логирование
│   └── utils/            # Утилиты (E.164 + тесты)
├── migrations/           # SQL миграции
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── sqlc.yaml
└── go.mod
📝 Makefile команды
Команда	Описание
make build	Сборка бинарника
make run	Запуск приложения
make docker-up	Запуск PostgreSQL
make generate	Генерация sqlc кода
make migrate-up	Применение миграций
make test	Запуск тестов
make clean	Очистка
🤝 Вклад в проект
PR принимаются. Для крупных изменений сначала откройте issue.

📄 Лицензия
MIT © kusyba

<div align="center">
⭐ Поставьте звезду, если проект вам полезен ⭐

</div>

