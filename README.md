# grpc-demo

[![CI](https://github.com/SergeiLPanov/grpc-demo/actions/workflows/ci.yml/badge.svg)](https://github.com/SergeiLPanov/grpc-demo/actions/workflows/ci.yml)

Учебный проект: два микросервиса на Go, общающиеся через gRPC. Запускается одной командой через Docker Compose.

## Архитектура

```
order-service  --(gRPC)-->  user-service
   (клиент)                    (сервер)
```

- **user-service** хранит пользователей и отдаёт их по gRPC (`GetUser`, `CreateUser`). Хранилище в памяти; благодаря интерфейсу `storage.UserStorage` его можно заменить на PostgreSQL, не меняя остальной код.
- **order-service** при создании заказа обращается к user-service и проверяет, что пользователь существует. Зависимость описана интерфейсом `UserChecker`, поэтому бизнес-логику можно тестировать без сети.
- **proto/users.proto** — общий контракт между сервисами.

## Что демонстрирует проект

- gRPC-сервер и клиент, контракт в `.proto`, сгенерированный код в `pb/`.
- Ошибки через gRPC-статусы: `NotFound`, `InvalidArgument`, `Internal`.
- Валидацию входных данных на стороне сервера.
- Таймаут на каждый вызов user-service и повторные попытки при старте.
- Корректную остановку сервера (`GracefulStop` по SIGTERM).
- Тесты: интеграционные через `bufconn` (настоящий gRPC в памяти) и юнит-тесты с заглушкой.
- Docker (multi-stage сборка), Docker Compose и CI на GitHub Actions.

## Структура

```
grpc-demo/
├── proto/users.proto
├── pb/                      # сгенерированный код
├── user-service/
│   ├── cmd/main.go
│   ├── internal/server/     # реализация gRPC-сервера и тесты
│   ├── internal/storage/    # хранилище в памяти
│   └── Dockerfile
├── order-service/
│   ├── cmd/main.go          # демо-сценарий
│   ├── internal/client/     # gRPC-клиент к user-service
│   ├── internal/service/    # бизнес-логика заказов и тесты
│   └── Dockerfile
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Запуск

### Через Docker Compose

```bash
docker compose up --build
```

order-service дождётся user-service, создаст пользователя, оформит ему заказ и попробует оформить заказ несуществующему пользователю. В логах появятся строки вида:

```
пользователь создан: <uuid>
заказ создан: {ID:<uuid> UserID:<uuid> Item:Клавиатура}
для несуществующего пользователя заказ не создан (так и должно быть): ...
```

После этого order-service завершится, а user-service продолжит работать на порту `50051`. Остановить всё: `Ctrl+C`, затем `docker compose down`.

### Локально, без Docker

Нужен Go 1.25 или новее.

```bash
go run ./user-service/cmd    # в одном терминале
go run ./order-service/cmd   # в другом
```

Адрес user-service для order-service задаётся переменной `USER_SERVICE_ADDR` (по умолчанию `localhost:50051`).

## Тесты

```bash
go test ./...
go test -race ./...
```

## Если меняете контракт

Сгенерированный код уже лежит в `pb/`. Если вы правите `proto/users.proto`, перегенерируйте его:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/users.proto
```

Нужен установленный `protoc`.

## Что можно добавить

- PostgreSQL вместо хранилища в памяти.
- REST-шлюз поверх order-service для внешних клиентов.
- Логирование и метрики (Prometheus).
- Health-check по gRPC и healthcheck в Docker Compose.
