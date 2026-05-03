ысоконагруженная система для отслеживания статистики GitHub-репозиториев в реальном времени.

## 🚀 Архитектура
Система построена на принципах **Clean Architecture** и **Event-Driven Design**:
- **API Gateway**: REST-вход, управление подписками и выдача закэшированных данных.
- **Processor**: "Мозг" системы. Хранит кэш в Postgres, управляет очередью задач в Kafka.
- **Collector**: Воркер, собирающий данные из GitHub REST API.
- **Subscriber**: gRPC-сервис управления списком подписок.

## 🛠 Стек технологий
- **Язык**: Go (Golang)
- **Очереди**: Apache Kafka
- **Базы данных**: PostgreSQL (2 независимых инстанса)
- **Интерфейсы**: gRPC, REST (Gin)
- **Инструменты**: SQLC, Golang-migrate, Docker Compose, Kafka-UI

## 🧩 Ключевые фичи
- **Idempotent Consumer (Inbox)**: Защита от дублей сообщений в Kafka через таблицу-инбокс по UUID.
- **Transactional Outbox**: Гарантированная отправка задач в очередь только после успешного коммита в БД.
- **Background Updates**: Автоматическое обновление данных по всем подпискам каждые 15 секунд через асинхронные задачи.
- **Parallel Processing**: Многопоточный сбор данных из GitHub (Worker Pool в Collector).

## 📦 Запуск
```bash
docker-compose up --build