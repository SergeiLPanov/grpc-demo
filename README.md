# grpc-demo

Учебный проект: два микросервиса на Go, общающиеся через gRPC.

## Архитектура

```
order-service  --(gRPC)-->  user-service
   (клиент)                    (сервер)
```

- **user-service** — хранит пользователей, отдаёт данные по gRPC. Хранилище — in-memory (легко заменить на PostgreSQL, реализовав интерфейс `storage.UserStorage`).
- **order-service** — при создании заказа обращается к user-service, чтобы проверить, что пользователь существует.
- **proto/users.proto** — общий контракт между сервисами.

## Структура

```
grpc-demo/
├── proto/users.proto
├── user-service/
│   ├── cmd/main.go
│   ├── internal/server/     # реализация gRPC-сервера
│   ├── internal/storage/    # in-memory хранилище
│   └── Dockerfile
├── order-service/
│   ├── cmd/main.go
│   ├── internal/client/     # gRPC-клиент к user-service
│   ├── internal/service/    # бизнес-логика заказов
│   └── Dockerfile
└── docker-compose.yml
```

## Как запустить

### 1. Сгенерировать Go-код из .proto (один раз)

Установите protoc и плагины:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Сгенерируйте код (создаст папку `pb/`):
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/users.proto
```

### 2. Установить зависимости

```bash
go mod init grpc-demo
go mod tidy
```

Понадобятся: `google.golang.org/grpc`, `google.golang.org/protobuf`, `github.com/google/uuid`, `github.com/stretchr/testify`.

### 3. Запустить локально (без Docker)

```bash
go run ./user-service/cmd    # в одном терминале
go run ./order-service/cmd   # в другом
```

### 4. Запустить через Docker Compose

```bash
docker compose up --build
```

## Тесты

```bash
go test ./... -v
go test -race ./...
```

## Что дальше можно добавить

- Заменить in-memory хранилище на PostgreSQL
- Добавить HTTP-обёртку (REST-шлюз) поверх order-service для внешних клиентов
- Добавить логирование и метрики (Prometheus)
