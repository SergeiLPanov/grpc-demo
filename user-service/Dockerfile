FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /user-service ./user-service/cmd

FROM alpine:latest
COPY --from=builder /user-service /user-service
EXPOSE 50051
ENTRYPOINT ["/user-service"]
