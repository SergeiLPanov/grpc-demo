FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /order-service ./order-service/cmd

FROM alpine:latest
COPY --from=builder /order-service /order-service
ENTRYPOINT ["/order-service"]
