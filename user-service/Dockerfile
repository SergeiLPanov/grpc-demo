FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /user-service ./user-service/cmd

FROM alpine:latest
COPY --from=builder /user-service /user-service
EXPOSE 50051
ENTRYPOINT ["/user-service"]
