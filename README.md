# RepoStat: Distributed Event-Driven Monitoring System

**RepoStat** — это высоконагруженная отказоустойчивая система для мониторинга статистики GitHub-репозиториев. Проект демонстрирует современные подходы к построению распределенных систем, асинхронную обработку данных и принципы Clean Architecture.

##  Архитектура системы

Система разделена на 4 независимых микросервиса:

- **API Gateway**: Входная точка (REST). Управляет подписками и предоставляет кэшированные данные. Также кэширует данные в Redis и ограничивает входящий трафик. 
- **Processor**: Оркестратор данных. Реализует логику кэширования в PostgreSQL, управляет транзакционными состояниями и очередями задач.
- **Collector**: Воркер-собиратель. Взаимодействует с GitHub REST API, обрабатывает задачи из Kafka и публикует результаты.
- **Subscriber**: Источник правды (Source of Truth) для управления списком пользовательских подписок.

### Data Flow (Поток данных)
`User` → `Gateway` → `gRPC` → `Subscriber` (Check subs)  
`Gateway` → `gRPC` → `Processor` (Cache check) → `Kafka (tasks)` → `Collector` → `GitHub API` → `Kafka (results)` → `Processor` (Update DB)

##  Технологический стек

- **Language**: Go 
- **Message Broker**: Apache Kafka 
- **Databases**: PostgreSQL 
- **Transport**: gRPC (Protobuf), REST 
- **Infrastructure**: Docker, Docker Compose, Golang-migrate

##  Ключевые инженерные фичи

###  Надежность (Exactly-Once Effect)
- **Transactional Outbox**: Гарантированная отправка задач в Kafka. Сообщение пишется в БД в одной транзакции с бизнес-логикой и доставляется отдельным Relay-воркером.
- **Transactional Inbox**: Идемпотентная обработка результатов. Каждый ответ из Kafka имеет уникальный UUID, который проверяется через таблицу-инбокс в Postgres перед обновлением кэша.

###  Оптимизация производительности
- **Cache Stampede Protection**: При первом запросе данных система ставит транзакционный "замок" (`status: FETCHING`), предотвращая дублирование задач в Kafka при наплыве пользователей.
- **Concurrent Processing**: Параллельное чтение партиций Kafka и многопоточный опрос GitHub API.

###  Синхронизация состояний
- **Background Polling**: Автоматическое обновление всех активных подписок каждые 15 секунд.
- **Event-Driven Invalidation**: Мгновенная очистка кэша в Processor при удалении подписки в Subscriber через шину событий.

### Caching in API Gateway(Redis):
- **Кэширование тяжелых gRPC-запросов (информация о репозиториях)**. Это позволило снизить нагрузку на сервис Processor и внешнее API GitHub, ускорив повторные ответы в несколько раз.

### Rate Limiting in API Gateway: 
- **Rate limiter**: реализован на базе Redis(алгоритм Fixed Window), также есть локальный rate limiter, который реализован с помощью алгоритма token Bucket, при деградации Redis, система переходит на локальный rate limiter.
- 
##  Запуск
## Запуск с помощью Docker
1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/Chaice1/RepoStat-Distributed-Async-Monitoring-System
    ```

2.  **Перейдите в директорию проекта:**
    ```bash
    cd RepoStat-Distributed-Async-Monitoring-System
    ```
3.  **Запуск программы:**
    ```bash
    make up
    ```
4.  **Остановка работы программы:**
    ```
    make down
    ```
5.  **Остановка работы программы и удаление volumes:**
    ```
    make down-volumes
    ```


## 🛠 API Endpoints

| Метод | Эндпоинт | Описание | Статусы |
| :--- | :--- | :--- | :--- |
| **GET** | `/api/ping`| проверить состояние системы. | `200`, `503`|
| **GET** | `/api/repositories/info?url=...` | Инфо о репозитории. Если нет в кэше — инициирует сбор. | `200`, `404`,`400`,`500` |
| **GET** | `/subscriptions` | Список всех активных подписок. | `200`, `500` |
| **POST** | `/subscriptions` | Подписка на репозиторий. Проверяет наличие на GitHub. | `200`, `400`, `404`,`500`,`409` |
| **DELETE** | `/subscriptions/{owner}/{repo}` | Отписка и очистка кэша в Processor. | `200`, `400`,`500` |
| **GET** | `/subscriptions/info` | Статистика по всем подпискам сразу (из кэша). | `200`,`500` |
