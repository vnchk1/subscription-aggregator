# Subscription Aggregator Service

## Описание сервиса

RESTful сервис для управления онлайн-подписками пользователей. Позволяет создавать, читать, обновлять, удалять и анализировать подписки с возможностью расчета общей стоимости за период.

### Основные функции:
- **CRUD операции** над подписками
- **Агрегация данных** - расчет общей стоимости подписок за период
- **Фильтрация** по пользователю и названию сервиса
- **Пагинация** для списка подписок
- **Swagger документация** для API

### Технологический стек:
- **Backend**: Go 1.24, Echo framework
- **База данных**: PostgreSQL 15
- **Миграции**: Goose
- **Контейнеризация**: Docker, Docker Compose
- **Документация**: Swagger/OpenAPI 3.0

### Минимальные требования

- **Git** (для клонирования репозитория)
- **Docker**: версия 20.10+
- **Docker Compose**: версия 2.0+
- **Оперативная память**: 2 GB минимум
- **Свободное место на диске**: 1 GB

## Быстрый старт

### 1. Клонирование и настройка
```bash
git clone https://github.com/vnchk1/subscription-aggregator.git <Ваша директория>
cd <Ваша директория>
```
### 2. Создание .env файла по примеру (для локальной разработки)
```bash
cp .env.example .env
```
** Либо настройка собственных переменных окружения
### 3. Запуск сервиса
```bash
docker-compose up --build
```

Сервис будет доступен по адресу: **http://localhost:8080**
### Для следующих запусков
```bash
docker-compose up
```

## Переменные окружения

Сервис использует `.env` файл со следующими настройками (для локальной разработки):

```env
# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10

# Database
DB_HOST=db
DB_PORT=5432
DB_NAME=sub_aggregator
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
DATABASE_MAX_CONNECTIONS=10
MIGRATION_PATH=./migrations

# Logger
LOG_LEVEL=debug
```

## API Endpoints

### Подписки
- `GET /subscriptions` - список подписок с пагинацией
- `POST /subscriptions` - создание подписки
- `GET /subscriptions/{id}` - получение подписки по ID
- `PUT /subscriptions/{id}` - обновление подписки
- `DELETE /subscriptions/{id}` - удаление подписки
- `GET /subscriptions/total-cost` - расчет общей стоимости

### Вспомогательные
- `GET /health` - health check
- `GET /swagger/index.html` - Swagger документация

### Остановка сервиса
```bash
docker-compose down
```

### Остановка сервиса, удаление контейнеров и томов
```bash
docker-compose down --rmi all --volumes --remove-orphans
```

## Мониторинг

- **Health check**: `http://localhost:8080/health`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **База данных**: `localhost:5432` (PostgreSQL)

### *Примечание для заказчика: .env файл не внесён в .gitignore, так как по ТЗ вся конфигурация должна храниться в нём*.