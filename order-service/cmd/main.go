package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-demo/order-service/internal/client"
	"grpc-demo/order-service/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr := os.Getenv("USER_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userClient, err := client.NewUserClient(addr)
	if err != nil {
		return err
	}
	defer userClient.Close()

	orderSvc := service.NewOrderService(userClient)

	// Сценарий демо: создаём пользователя в user-service,
	// оформляем ему заказ, затем пробуем оформить заказ несуществующему.
	userID, err := createUserWithRetry(ctx, userClient)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}
	log.Printf("пользователь создан: %s", userID)

	order, err := orderSvc.CreateOrder(ctx, userID, "Клавиатура")
	if err != nil {
		return fmt.Errorf("не удалось создать заказ: %w", err)
	}
	log.Printf("заказ создан: %+v", order)

	_, err = orderSvc.CreateOrder(ctx, "unknown-user-id", "Мышь")
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		log.Printf("для несуществующего пользователя заказ не создан (так и должно быть): %v", err)
	case err != nil:
		return fmt.Errorf("неожиданная ошибка: %w", err)
	default:
		return errors.New("заказ для несуществующего пользователя не должен создаваться")
	}

	return nil
}

// createUserWithRetry ждёт, пока user-service запустится: в docker compose
// order-service может стартовать раньше, чем сервер начнёт принимать запросы.
func createUserWithRetry(ctx context.Context, c *client.UserClient) (string, error) {
	const attempts = 10

	var lastErr error
	for i := 1; i <= attempts; i++ {
		id, err := c.CreateUser(ctx, "Sergei", "sergei@example.com")
		if err == nil {
			return id, nil
		}
		lastErr = err

		code := status.Code(err)
		if code != codes.Unavailable && code != codes.DeadlineExceeded {
			return "", err
		}

		log.Printf("user-service пока недоступен (попытка %d/%d), жду...", i, attempts)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return "", fmt.Errorf("user-service недоступен: %w", lastErr)
}
