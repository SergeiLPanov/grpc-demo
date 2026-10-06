FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /order-service ./order-service/cmd

FROM alpine:latest
COPY --from=builder /order-service /order-service
ENTRYPOINT ["/order-service"]
