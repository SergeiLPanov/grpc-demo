package service

import (
	"context"
	"errors"
	"sync"

	"grpc-demo/order-service/internal/client"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("пользователь не найден в user-service")

type Order struct {
	ID     string
	UserID string
	Item   string
}

type OrderService struct {
	mu         sync.RWMutex
	orders     map[string]Order
	userClient *client.UserClient
}

func NewOrderService(userClient *client.UserClient) *OrderService {
	return &OrderService{
		orders:     make(map[string]Order),
		userClient: userClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, item string) (Order, error) {
	exists, err := s.userClient.UserExists(ctx, userID)
	if err != nil {
		return Order{}, err
	}
	if !exists {
		return Order{}, ErrUserNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order := Order{
		ID:     uuid.NewString(),
		UserID: userID,
		Item:   item,
	}
	s.orders[order.ID] = order
	return order, nil
}
